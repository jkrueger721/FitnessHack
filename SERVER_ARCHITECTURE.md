# Fitness Hack Server Architecture

## Overview

The Fitness Hack server is a Go-based REST API built with modern architectural patterns, featuring clean separation of concerns, comprehensive caching, authentication, and robust error handling.

## Architecture Overview

```
┌─────────────────┐    ┌─────────────────┐    ┌─────────────────┐
│   HTTP Client   │───▶│   Fiber Server  │───▶│   Database      │
│                 │    │                 │    │   (PostgreSQL)  │
└─────────────────┘    └─────────────────┘    └─────────────────┘
                              │
                              ▼
                       ┌─────────────────┐
                       │   Redis Cache   │
                       └─────────────────┘
```

## Directory Structure

```
internal/server/
├── server.go           # Main server setup and configuration
├── routes.go           # Route definitions and middleware setup
├── utils.go            # Shared utility functions
├── users.go            # User-related handlers
├── workouts.go         # Workout-related handlers
├── exercises.go        # Exercise-related handlers
├── workout_exercises.go # Workout-exercise relationship handlers
├── workout_sessions.go # Workout session handlers
└── routes_test.go      # Route testing utilities
```

## Core Components

### 1. Server Setup (`server.go`)

The main server file handles:
- **Fiber App Configuration**: Middleware setup, error handling, logging
- **Database Connection**: PostgreSQL connection with connection pooling
- **Redis Cache**: Cache client initialization and configuration
- **Graceful Shutdown**: Proper cleanup on server termination

#### Key Features:
- **Connection Pooling**: Configurable database connection pool
- **Structured Logging**: JSON-formatted logs optimized for CloudWatch
- **Error Handling**: Centralized error handling with stack traces
- **Health Checks**: Database and cache health monitoring

### 2. Route Management (`routes.go`)

Centralized route configuration with:
- **API Versioning**: `/api/v1` prefix for all routes
- **Middleware Stack**: Authentication, logging, CORS, rate limiting
- **Resource Grouping**: Logical grouping of related endpoints
- **Error Recovery**: Panic recovery middleware

#### Route Structure:
```go
/api/v1/
├── auth/
│   └── login          # POST - User authentication
├── users/             # User management
├── workouts/          # Workout management
├── exercises/         # Exercise management
├── workout-exercises/ # Workout-exercise relationships
└── workout-sessions/  # Workout session tracking
```

### 3. Handler Architecture

Each resource has its own handler file following consistent patterns:

#### Handler Structure:
```go
// Cache key helpers
func resourceCacheKey(id string) string
func resourcesListCacheKey(limit, offset int) string

// Model conversion helpers
func resourceToResponse(resource *database.Resource) database.ResourceResponse

// CRUD handlers
func (s *FiberServer) createResource(c *fiber.Ctx) error
func (s *FiberServer) getResource(c *fiber.Ctx) error
func (s *FiberServer) listResources(c *fiber.Ctx) error
func (s *FiberServer) updateResource(c *fiber.Ctx) error
func (s *FiberServer) deleteResource(c *fiber.Ctx) error
```

## Design Patterns

### 1. Repository Pattern

The database layer implements the repository pattern through the `Service` interface:

```go
type Service interface {
    // Health and connection management
    Health() map[string]string
    Close() error
    GetDB() *sqlx.DB
    BeginTx(ctx context.Context) (*sqlx.Tx, error)
    
    // CRUD operations for each resource
    CreateUser(ctx context.Context, user *Users) (*Users, error)
    GetUserByID(ctx context.Context, id string) (*Users, error)
    // ... other methods
}
```

### 2. Request/Response Models

Clean separation between database models and API models:

#### Database Models:
- Represent the actual database schema
- Include all fields including sensitive data
- Used for internal operations

#### Request Models:
- Validate and sanitize input data
- Include validation tags
- Handle optional fields with pointers

#### Response Models:
- Exclude sensitive data (e.g., password hashes)
- Optimized for API consumption
- Include computed fields if needed

### 3. Caching Strategy

Multi-level caching approach:

#### Cache Levels:
1. **Individual Resource Cache**: `{resource}:{id}`
2. **List Resource Cache**: `{resource}:list:{limit}:{offset}`
3. **Cache Duration**: 10 minutes for most resources

#### Cache Operations:
- **Read**: Check cache first, fallback to database
- **Write**: Update database, invalidate related caches
- **Delete**: Remove from database, invalidate caches

#### Cache Invalidation:
```go
// Invalidate individual cache
s.DeleteCache(ctx, resourceCacheKey(id))

// Invalidate list caches
s.cache.Del(ctx, "resources:list:*")
```

### 4. Authentication & Authorization

JWT-based authentication with middleware:

#### JWT Middleware:
- Extracts token from Authorization header
- Validates token signature and expiration
- Sets user context for downstream handlers

#### Protected Routes:
- All routes except `/auth/login` and `/users` (registration)
- User ID automatically available in handlers via `c.Locals("user_id")`

#### Security Features:
- Password hashing with bcrypt
- JWT token expiration (24 hours)
- Secure password validation

## Error Handling

### 1. Structured Error Responses

Consistent error format across all endpoints:

```json
{
  "error": {
    "message": "Human-readable error message",
    "code": "ERROR_CODE",
    "details": {}
  }
}
```

### 2. Error Types

- **Validation Errors**: Request data validation failures
- **Authentication Errors**: Invalid or missing JWT tokens
- **Authorization Errors**: Insufficient permissions
- **Not Found Errors**: Requested resource doesn't exist
- **Internal Errors**: Unexpected server errors

### 3. Error Logging

Structured error logging with:
- **Stack Traces**: Full error context for debugging
- **Request Context**: User ID, request path, method
- **Error Classification**: Different log levels for different error types

## Database Layer

### 1. Connection Management

```go
type Config struct {
    MaxOpenConns    int
    MaxIdleConns    int
    ConnMaxLifetime time.Duration
    ConnMaxIdleTime time.Duration
}
```

### 2. Transaction Support

```go
func (s *service) BeginTx(ctx context.Context) (*sqlx.Tx, error)
```

### 3. Health Monitoring

```go
func (s *service) Health() map[string]string
```

## Caching Layer

### 1. Redis Configuration

- **Connection Pooling**: Multiple connections for concurrent access
- **Key Expiration**: Automatic cleanup of expired keys
- **Pattern Matching**: Support for wildcard cache invalidation

### 2. Cache Operations

```go
// Set cache with expiration
s.SetCache(ctx, key, value, 10*time.Minute)

// Get cache value
value, err := s.GetCache(ctx, key)

// Delete cache
s.DeleteCache(ctx, key)

// Pattern-based deletion
s.cache.Del(ctx, "pattern:*")
```

## Middleware Stack

### 1. Request Logging

Structured request logging with:
- **Request ID**: Unique identifier for request tracing
- **Method & Path**: HTTP method and request path
- **Response Time**: Request processing duration
- **Status Code**: HTTP response status

### 2. CORS Configuration

```go
app.Use(cors.New(cors.Config{
    AllowOrigins: "*",
    AllowMethods: "GET,POST,PUT,DELETE,OPTIONS",
    AllowHeaders: "Origin,Content-Type,Accept,Authorization",
}))
```

### 3. Rate Limiting

Currently not implemented but can be easily added using Fiber's rate limiting middleware.

### 4. Panic Recovery

```go
app.Use(recover.New(recover.Config{
    EnableStackTrace: true,
}))
```

## Testing Strategy

### 1. Unit Tests

- **Handler Tests**: Test individual handler functions
- **Model Tests**: Test request/response model validation
- **Utility Tests**: Test helper functions

### 2. Integration Tests

- **Database Tests**: Test database operations with test containers
- **Cache Tests**: Test Redis operations
- **API Tests**: Test complete request/response cycles

### 3. Test Utilities

```go
// Test database setup
func setupTestDB() (*sqlx.DB, func())

// Test cache setup
func setupTestCache() (*redis.Client, func())

// Test server setup
func setupTestServer() (*FiberServer, func())
```

## Performance Considerations

### 1. Database Optimization

- **Connection Pooling**: Reuse database connections
- **Query Optimization**: Use prepared statements
- **Indexing**: Proper database indexes for common queries

### 2. Caching Strategy

- **Cache Hit Ratio**: Monitor cache effectiveness
- **Cache Warming**: Pre-populate frequently accessed data
- **Cache Size**: Monitor memory usage

### 3. Response Optimization

- **JSON Marshaling**: Efficient JSON serialization
- **Response Compression**: Gzip compression for large responses
- **Pagination**: Limit response size with pagination

## Security Considerations

### 1. Input Validation

- **Request Validation**: Validate all input data
- **SQL Injection Prevention**: Use parameterized queries
- **XSS Prevention**: Sanitize user input

### 2. Authentication Security

- **Password Hashing**: bcrypt with appropriate cost
- **JWT Security**: Secure token generation and validation
- **Token Expiration**: Automatic token expiration

### 3. Data Protection

- **Sensitive Data**: Never return sensitive data in responses
- **Data Encryption**: Encrypt sensitive data at rest
- **Access Control**: Proper authorization checks

## Monitoring & Observability

### 1. Health Checks

```go
func (s *service) Health() map[string]string {
    return map[string]string{
        "status": "up",
        "database": "connected",
        "cache": "connected",
        "connections": "25/25",
    }
}
```

### 2. Metrics Collection

- **Request Count**: Number of requests per endpoint
- **Response Time**: Average response time
- **Error Rate**: Percentage of failed requests
- **Cache Hit Ratio**: Cache effectiveness

### 3. Logging

- **Structured Logs**: JSON-formatted logs
- **Log Levels**: Different levels for different types of events
- **Request Tracing**: Track requests across the system

## Deployment Considerations

### 1. Environment Configuration

```bash
# Database
BLUEPRINT_DB_HOST=localhost
BLUEPRINT_DB_PORT=5432
BLUEPRINT_DB_DATABASE=fitness_hack
BLUEPRINT_DB_USERNAME=postgres
BLUEPRINT_DB_PASSWORD=password
BLUEPRINT_DB_SCHEMA=public

# Redis
REDIS_HOST=localhost
REDIS_PORT=6379
REDIS_PASSWORD=
REDIS_DB=0

# JWT
JWT_SECRET=your-secret-key

# Server
PORT=8080
ENV=development
```

### 2. Docker Support

```dockerfile
FROM golang:1.21-alpine
WORKDIR /app
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN go build -o main .
EXPOSE 8080
CMD ["./main"]
```

### 3. Kubernetes Deployment

- **Horizontal Pod Autoscaling**: Scale based on CPU/memory usage
- **Health Checks**: Liveness and readiness probes
- **Resource Limits**: CPU and memory limits
- **Secrets Management**: Secure configuration management

## Future Enhancements

### 1. Planned Features

- **Rate Limiting**: Implement request rate limiting
- **API Versioning**: Support for multiple API versions
- **GraphQL**: Add GraphQL endpoint
- **WebSocket**: Real-time updates for workout sessions

### 2. Performance Improvements

- **Database Sharding**: Horizontal database scaling
- **CDN Integration**: Static asset delivery
- **Microservices**: Split into smaller services
- **Event Sourcing**: Event-driven architecture

### 3. Security Enhancements

- **OAuth Integration**: Social login support
- **2FA**: Two-factor authentication
- **API Keys**: API key management
- **Audit Logging**: Comprehensive audit trail

## Conclusion

The Fitness Hack server follows modern architectural patterns with a focus on:
- **Scalability**: Horizontal scaling capabilities
- **Maintainability**: Clean code structure and separation of concerns
- **Security**: Comprehensive security measures
- **Performance**: Efficient caching and database operations
- **Observability**: Comprehensive monitoring and logging

The architecture provides a solid foundation for future enhancements while maintaining high performance and reliability standards. 