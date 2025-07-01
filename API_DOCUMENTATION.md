# Fitness Hack API Documentation

## Overview

The Fitness Hack API is a RESTful service built with Go, Fiber, and PostgreSQL that provides comprehensive fitness tracking capabilities. The API supports user management, workout planning, exercise tracking, and session logging with JWT authentication and Redis caching.

## Table of Contents

- [Authentication](#authentication)
- [Base URL](#base-url)
- [Error Handling](#error-handling)
- [Pagination](#pagination)
- [Endpoints](#endpoints)
  - [Authentication](#authentication-endpoints)
  - [Users](#users-endpoints)
  - [Workouts](#workouts-endpoints)
  - [Exercises](#exercises-endpoints)
  - [Workout Exercises](#workout-exercises-endpoints)
  - [Workout Sessions](#workout-sessions-endpoints)
- [Data Models](#data-models)
- [Caching](#caching)
- [Rate Limiting](#rate-limiting)

## Authentication

The API uses JWT (JSON Web Tokens) for authentication. Most endpoints require a valid JWT token in the Authorization header.

### JWT Token Format
```
Authorization: Bearer <your-jwt-token>
```

### Token Expiration
- Tokens expire after 24 hours
- Refresh tokens are not currently supported

## Base URL

```
http://localhost:8080/api/v1
```

## Error Handling

All API responses follow a consistent error format:

### Error Response Format
```json
{
  "error": {
    "message": "Error description",
    "code": "ERROR_CODE",
    "details": {}
  }
}
```

### Common HTTP Status Codes
- `200` - Success
- `201` - Created
- `400` - Bad Request
- `401` - Unauthorized
- `404` - Not Found
- `500` - Internal Server Error

## Pagination

List endpoints support pagination with the following query parameters:

- `limit` (int, 1-100): Number of items per page (default: 10)
- `offset` (int, 0+): Number of items to skip (default: 0)

### Paginated Response Format
```json
{
  "data": [...],
  "pagination": {
    "limit": 10,
    "offset": 0,
    "total": 100
  }
}
```

## Endpoints

### Authentication Endpoints

#### POST /auth/login
Authenticate a user and receive a JWT token.

**Request Body:**
```json
{
  "email": "user@example.com",
  "password": "password123"
}
```

**Response:**
```json
{
  "data": {
    "user": {
      "id": "uuid",
      "email": "user@example.com",
      "username": "username",
      "first_name": "John",
      "last_name": "Doe",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    },
    "token": "jwt-token-here"
  }
}
```

### Users Endpoints

#### POST /users
Create a new user account.

**Request Body:**
```json
{
  "email": "user@example.com",
  "username": "username",
  "password": "password123",
  "first_name": "John",
  "last_name": "Doe"
}
```

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "username",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### GET /users/{id}
Get a specific user by ID.

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "email": "user@example.com",
    "username": "username",
    "first_name": "John",
    "last_name": "Doe",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### GET /users
Get a paginated list of users.

**Query Parameters:**
- `limit` (optional): Number of users per page
- `offset` (optional): Number of users to skip

**Response:**
```json
{
  "data": [
    {
      "id": "uuid",
      "email": "user@example.com",
      "username": "username",
      "first_name": "John",
      "last_name": "Doe",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "limit": 10,
    "offset": 0,
    "total": 1
  }
}
```

#### PUT /users/{id}
Update a user's information.

**Request Body:**
```json
{
  "email": "newemail@example.com",
  "username": "newusername",
  "password": "newpassword123",
  "first_name": "Jane",
  "last_name": "Smith"
}
```

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "email": "newemail@example.com",
    "username": "newusername",
    "first_name": "Jane",
    "last_name": "Smith",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

#### DELETE /users/{id}
Delete a user account.

**Response:** `204 No Content`

### Workouts Endpoints

#### POST /workouts
Create a new workout plan.

**Headers:** `Authorization: Bearer <jwt-token>`

**Request Body:**
```json
{
  "name": "Upper Body Strength",
  "description": "Focus on chest, back, and arms",
  "duration_minutes": 60
}
```

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "user_id": "user-uuid",
    "name": "Upper Body Strength",
    "description": "Focus on chest, back, and arms",
    "duration_minutes": 60,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### GET /workouts/{id}
Get a specific workout by ID.

**Headers:** `Authorization: Bearer <jwt-token>`

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "user_id": "user-uuid",
    "name": "Upper Body Strength",
    "description": "Focus on chest, back, and arms",
    "duration_minutes": 60,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### GET /workouts
Get a paginated list of workouts.

**Headers:** `Authorization: Bearer <jwt-token>`

**Query Parameters:**
- `limit` (optional): Number of workouts per page
- `offset` (optional): Number of workouts to skip

**Response:**
```json
{
  "data": [
    {
      "id": "uuid",
      "user_id": "user-uuid",
      "name": "Upper Body Strength",
      "description": "Focus on chest, back, and arms",
      "duration_minutes": 60,
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "limit": 10,
    "offset": 0,
    "total": 1
  }
}
```

#### PUT /workouts/{id}
Update a workout plan.

**Headers:** `Authorization: Bearer <jwt-token>`

**Request Body:**
```json
{
  "name": "Updated Upper Body",
  "description": "Updated description",
  "duration_minutes": 75
}
```

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "user_id": "user-uuid",
    "name": "Updated Upper Body",
    "description": "Updated description",
    "duration_minutes": 75,
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

#### DELETE /workouts/{id}
Delete a workout plan.

**Headers:** `Authorization: Bearer <jwt-token>`

**Response:** `204 No Content`

### Exercises Endpoints

#### POST /exercises
Create a new exercise.

**Request Body:**
```json
{
  "name": "Bench Press",
  "description": "Compound chest exercise",
  "muscle_group": "Chest",
  "equipment": "Barbell",
  "difficulty_level": "Intermediate",
  "instructions": "Lie on bench, lower bar to chest, press up"
}
```

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "name": "Bench Press",
    "description": "Compound chest exercise",
    "muscle_group": "Chest",
    "equipment": "Barbell",
    "difficulty_level": "Intermediate",
    "instructions": "Lie on bench, lower bar to chest, press up",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### GET /exercises/{id}
Get a specific exercise by ID.

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "name": "Bench Press",
    "description": "Compound chest exercise",
    "muscle_group": "Chest",
    "equipment": "Barbell",
    "difficulty_level": "Intermediate",
    "instructions": "Lie on bench, lower bar to chest, press up",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T00:00:00Z"
  }
}
```

#### GET /exercises
Get a paginated list of exercises.

**Query Parameters:**
- `limit` (optional): Number of exercises per page
- `offset` (optional): Number of exercises to skip

**Response:**
```json
{
  "data": [
    {
      "id": "uuid",
      "name": "Bench Press",
      "description": "Compound chest exercise",
      "muscle_group": "Chest",
      "equipment": "Barbell",
      "difficulty_level": "Intermediate",
      "instructions": "Lie on bench, lower bar to chest, press up",
      "created_at": "2024-01-01T00:00:00Z",
      "updated_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "limit": 10,
    "offset": 0,
    "total": 1
  }
}
```

#### PUT /exercises/{id}
Update an exercise.

**Request Body:**
```json
{
  "name": "Updated Bench Press",
  "description": "Updated description",
  "muscle_group": "Chest, Triceps",
  "equipment": "Barbell, Bench",
  "difficulty_level": "Advanced",
  "instructions": "Updated instructions"
}
```

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "name": "Updated Bench Press",
    "description": "Updated description",
    "muscle_group": "Chest, Triceps",
    "equipment": "Barbell, Bench",
    "difficulty_level": "Advanced",
    "instructions": "Updated instructions",
    "created_at": "2024-01-01T00:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

#### DELETE /exercises/{id}
Delete an exercise.

**Response:** `204 No Content`

### Workout Exercises Endpoints

#### POST /workout-exercises
Add an exercise to a workout.

**Request Body:**
```json
{
  "workout_id": "workout-uuid",
  "exercise_id": "exercise-uuid",
  "sets": 3,
  "reps": 10,
  "weight_kg": "100.5",
  "duration_seconds": 60,
  "order_index": 1,
  "rest_seconds": 90,
  "notes": "Focus on form"
}
```

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "workout_id": "workout-uuid",
    "exercise_id": "exercise-uuid",
    "sets": 3,
    "reps": 10,
    "weight_kg": "100.5",
    "duration_seconds": 60,
    "order_index": 1,
    "rest_seconds": 90,
    "notes": "Focus on form",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### GET /workout-exercises/{id}
Get a specific workout exercise by ID.

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "workout_id": "workout-uuid",
    "exercise_id": "exercise-uuid",
    "sets": 3,
    "reps": 10,
    "weight_kg": "100.5",
    "duration_seconds": 60,
    "order_index": 1,
    "rest_seconds": 90,
    "notes": "Focus on form",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### GET /workout-exercises
Get a paginated list of workout exercises.

**Query Parameters:**
- `limit` (optional): Number of workout exercises per page
- `offset` (optional): Number of workout exercises to skip

**Response:**
```json
{
  "data": [
    {
      "id": "uuid",
      "workout_id": "workout-uuid",
      "exercise_id": "exercise-uuid",
      "sets": 3,
      "reps": 10,
      "weight_kg": "100.5",
      "duration_seconds": 60,
      "order_index": 1,
      "rest_seconds": 90,
      "notes": "Focus on form",
      "created_at": "2024-01-01T00:00:00Z"
    }
  ],
  "pagination": {
    "limit": 10,
    "offset": 0,
    "total": 1
  }
}
```

#### PUT /workout-exercises/{id}
Update a workout exercise.

**Request Body:**
```json
{
  "sets": 4,
  "reps": 12,
  "weight_kg": "110.0",
  "rest_seconds": 120,
  "notes": "Updated notes"
}
```

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "workout_id": "workout-uuid",
    "exercise_id": "exercise-uuid",
    "sets": 4,
    "reps": 12,
    "weight_kg": "110.0",
    "duration_seconds": 60,
    "order_index": 1,
    "rest_seconds": 120,
    "notes": "Updated notes",
    "created_at": "2024-01-01T00:00:00Z"
  }
}
```

#### DELETE /workout-exercises/{id}
Remove an exercise from a workout.

**Response:** `204 No Content`

### Workout Sessions Endpoints

#### POST /workout-sessions
Start a new workout session.

**Headers:** `Authorization: Bearer <jwt-token>`

**Request Body:**
```json
{
  "workout_id": "workout-uuid",
  "name": "Morning Workout",
  "started_at": "2024-01-01T08:00:00Z",
  "completed_at": "2024-01-01T09:00:00Z",
  "duration_minutes": 60,
  "notes": "Great session!"
}
```

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "user_id": "user-uuid",
    "workout_id": "workout-uuid",
    "name": "Morning Workout",
    "started_at": "2024-01-01T08:00:00Z",
    "completed_at": "2024-01-01T09:00:00Z",
    "duration_minutes": 60,
    "notes": "Great session!",
    "created_at": "2024-01-01T08:00:00Z",
    "updated_at": "2024-01-01T08:00:00Z"
  }
}
```

#### GET /workout-sessions/{id}
Get a specific workout session by ID.

**Headers:** `Authorization: Bearer <jwt-token>`

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "user_id": "user-uuid",
    "workout_id": "workout-uuid",
    "name": "Morning Workout",
    "started_at": "2024-01-01T08:00:00Z",
    "completed_at": "2024-01-01T09:00:00Z",
    "duration_minutes": 60,
    "notes": "Great session!",
    "created_at": "2024-01-01T08:00:00Z",
    "updated_at": "2024-01-01T08:00:00Z"
  }
}
```

#### GET /workout-sessions
Get a paginated list of workout sessions.

**Headers:** `Authorization: Bearer <jwt-token>`

**Query Parameters:**
- `limit` (optional): Number of workout sessions per page
- `offset` (optional): Number of workout sessions to skip

**Response:**
```json
{
  "data": [
    {
      "id": "uuid",
      "user_id": "user-uuid",
      "workout_id": "workout-uuid",
      "name": "Morning Workout",
      "started_at": "2024-01-01T08:00:00Z",
      "completed_at": "2024-01-01T09:00:00Z",
      "duration_minutes": 60,
      "notes": "Great session!",
      "created_at": "2024-01-01T08:00:00Z",
      "updated_at": "2024-01-01T08:00:00Z"
    }
  ],
  "pagination": {
    "limit": 10,
    "offset": 0,
    "total": 1
  }
}
```

#### PUT /workout-sessions/{id}
Update a workout session.

**Headers:** `Authorization: Bearer <jwt-token>`

**Request Body:**
```json
{
  "completed_at": "2024-01-01T09:15:00Z",
  "duration_minutes": 75,
  "notes": "Updated notes"
}
```

**Response:**
```json
{
  "data": {
    "id": "uuid",
    "user_id": "user-uuid",
    "workout_id": "workout-uuid",
    "name": "Morning Workout",
    "started_at": "2024-01-01T08:00:00Z",
    "completed_at": "2024-01-01T09:15:00Z",
    "duration_minutes": 75,
    "notes": "Updated notes",
    "created_at": "2024-01-01T08:00:00Z",
    "updated_at": "2024-01-01T12:00:00Z"
  }
}
```

#### DELETE /workout-sessions/{id}
Delete a workout session.

**Headers:** `Authorization: Bearer <jwt-token>`

**Response:** `204 No Content`

## Data Models

### User Models

#### CreateUserRequest
```json
{
  "email": "string (required, email format)",
  "username": "string (required, 3-100 chars)",
  "password": "string (required, min 8 chars)",
  "first_name": "string (optional)",
  "last_name": "string (optional)"
}
```

#### UpdateUserRequest
```json
{
  "email": "string (optional, email format)",
  "username": "string (optional, 3-100 chars)",
  "password": "string (optional, min 8 chars)",
  "first_name": "string (optional)",
  "last_name": "string (optional)"
}
```

#### UserResponse
```json
{
  "id": "string (UUID)",
  "email": "string",
  "username": "string",
  "first_name": "string (optional)",
  "last_name": "string (optional)",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Workout Models

#### CreateWorkoutRequest
```json
{
  "name": "string (required, max 255 chars)",
  "description": "string (optional)",
  "duration_minutes": "integer (optional, min 1)"
}
```

#### UpdateWorkoutRequest
```json
{
  "name": "string (optional, max 255 chars)",
  "description": "string (optional)",
  "duration_minutes": "integer (optional, min 1)"
}
```

#### WorkoutResponse
```json
{
  "id": "string (UUID)",
  "user_id": "string (UUID)",
  "name": "string",
  "description": "string (optional)",
  "duration_minutes": "integer (optional)",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Exercise Models

#### CreateExerciseRequest
```json
{
  "name": "string (required, max 255 chars)",
  "description": "string (optional)",
  "muscle_group": "string (optional, max 100 chars)",
  "equipment": "string (optional, max 100 chars)",
  "difficulty_level": "string (optional, max 50 chars)",
  "instructions": "string (optional)"
}
```

#### UpdateExerciseRequest
```json
{
  "name": "string (optional, max 255 chars)",
  "description": "string (optional)",
  "muscle_group": "string (optional, max 100 chars)",
  "equipment": "string (optional, max 100 chars)",
  "difficulty_level": "string (optional, max 50 chars)",
  "instructions": "string (optional)"
}
```

#### ExerciseResponse
```json
{
  "id": "string (UUID)",
  "name": "string",
  "description": "string (optional)",
  "muscle_group": "string (optional)",
  "equipment": "string (optional)",
  "difficulty_level": "string (optional)",
  "instructions": "string (optional)",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

### Workout Exercise Models

#### CreateWorkoutExerciseRequest
```json
{
  "workout_id": "string (required, UUID)",
  "exercise_id": "string (required, UUID)",
  "sets": "integer (required, min 1)",
  "reps": "integer (optional, min 1)",
  "weight_kg": "decimal (optional, min 0)",
  "duration_seconds": "integer (optional, min 1)",
  "order_index": "integer (min 0)",
  "rest_seconds": "integer (min 0)",
  "notes": "string (optional)"
}
```

#### UpdateWorkoutExerciseRequest
```json
{
  "workout_id": "string (optional, UUID)",
  "exercise_id": "string (optional, UUID)",
  "sets": "integer (optional, min 1)",
  "reps": "integer (optional, min 1)",
  "weight_kg": "decimal (optional, min 0)",
  "duration_seconds": "integer (optional, min 1)",
  "order_index": "integer (optional, min 0)",
  "rest_seconds": "integer (optional, min 0)",
  "notes": "string (optional)"
}
```

#### WorkoutExerciseResponse
```json
{
  "id": "string (UUID)",
  "workout_id": "string (UUID)",
  "exercise_id": "string (UUID)",
  "sets": "integer",
  "reps": "integer (optional)",
  "weight_kg": "decimal (optional)",
  "duration_seconds": "integer (optional)",
  "order_index": "integer",
  "rest_seconds": "integer",
  "notes": "string (optional)",
  "created_at": "datetime"
}
```

### Workout Session Models

#### CreateWorkoutSessionRequest
```json
{
  "workout_id": "string (optional, UUID)",
  "name": "string (required, max 255 chars)",
  "started_at": "datetime (optional)",
  "completed_at": "datetime (optional)",
  "duration_minutes": "integer (optional, min 1)",
  "notes": "string (optional)"
}
```

#### UpdateWorkoutSessionRequest
```json
{
  "workout_id": "string (optional, UUID)",
  "name": "string (optional, max 255 chars)",
  "started_at": "datetime (optional)",
  "completed_at": "datetime (optional)",
  "duration_minutes": "integer (optional, min 1)",
  "notes": "string (optional)"
}
```

#### WorkoutSessionResponse
```json
{
  "id": "string (UUID)",
  "user_id": "string (UUID)",
  "workout_id": "string (optional, UUID)",
  "name": "string",
  "started_at": "datetime",
  "completed_at": "datetime (optional)",
  "duration_minutes": "integer (optional)",
  "notes": "string (optional)",
  "created_at": "datetime",
  "updated_at": "datetime"
}
```

## Caching

The API uses Redis for caching to improve performance:

### Cache Keys
- Individual resources: `{resource}:{id}` (e.g., `user:uuid`)
- List resources: `{resource}:list:{limit}:{offset}` (e.g., `users:list:10:0`)
- Cache duration: 10 minutes for most resources

### Cache Invalidation
- Cache is automatically invalidated on create, update, and delete operations
- List caches are invalidated when individual resources are modified

## Rate Limiting

Currently, the API does not implement rate limiting. Consider implementing rate limiting for production use.

## Usage Examples

### Complete Workflow Example

1. **Create a user account:**
```bash
curl -X POST http://localhost:8080/api/v1/users \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "username": "fitnessuser",
    "password": "password123",
    "first_name": "John",
    "last_name": "Doe"
  }'
```

2. **Login to get JWT token:**
```bash
curl -X POST http://localhost:8080/api/v1/auth/login \
  -H "Content-Type: application/json" \
  -d '{
    "email": "user@example.com",
    "password": "password123"
  }'
```

3. **Create a workout:**
```bash
curl -X POST http://localhost:8080/api/v1/workouts \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt-token>" \
  -d '{
    "name": "Upper Body Strength",
    "description": "Focus on chest, back, and arms",
    "duration_minutes": 60
  }'
```

4. **Create exercises:**
```bash
curl -X POST http://localhost:8080/api/v1/exercises \
  -H "Content-Type: application/json" \
  -d '{
    "name": "Bench Press",
    "description": "Compound chest exercise",
    "muscle_group": "Chest",
    "equipment": "Barbell",
    "difficulty_level": "Intermediate",
    "instructions": "Lie on bench, lower bar to chest, press up"
  }'
```

5. **Add exercise to workout:**
```bash
curl -X POST http://localhost:8080/api/v1/workout-exercises \
  -H "Content-Type: application/json" \
  -d '{
    "workout_id": "workout-uuid",
    "exercise_id": "exercise-uuid",
    "sets": 3,
    "reps": 10,
    "weight_kg": "100.5",
    "order_index": 1,
    "rest_seconds": 90
  }'
```

6. **Start a workout session:**
```bash
curl -X POST http://localhost:8080/api/v1/workout-sessions \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt-token>" \
  -d '{
    "workout_id": "workout-uuid",
    "name": "Morning Workout",
    "notes": "Starting my workout"
  }'
```

7. **Complete the workout session:**
```bash
curl -X PUT http://localhost:8080/api/v1/workout-sessions/session-uuid \
  -H "Content-Type: application/json" \
  -H "Authorization: Bearer <jwt-token>" \
  -d '{
    "completed_at": "2024-01-01T09:00:00Z",
    "duration_minutes": 60,
    "notes": "Great workout! Felt strong today."
  }'
```

## Error Codes

| Code | Description |
|------|-------------|
| `INVALID_REQUEST` | Request body is malformed or missing required fields |
| `UNAUTHORIZED` | Invalid or missing JWT token |
| `FORBIDDEN` | User doesn't have permission to access resource |
| `NOT_FOUND` | Requested resource doesn't exist |
| `VALIDATION_ERROR` | Request data fails validation rules |
| `INTERNAL_ERROR` | Server encountered an unexpected error |

## Support

For questions or issues, please refer to the project repository or contact the development team. 