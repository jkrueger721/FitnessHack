package server

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"time"

	"fitness-hack/internal/database"

	"github.com/gofiber/fiber/v2"
	"github.com/golang-jwt/jwt/v5"
	"golang.org/x/crypto/bcrypt"
)

// Cache key helpers
func userCacheKey(id string) string {
	return fmt.Sprintf("user:%s", id)
}

func usersListCacheKey(limit, offset int) string {
	return fmt.Sprintf("users:list:%d:%d", limit, offset)
}

// Helper to hash password
func hashPassword(password string) (string, error) {
	bytes, err := bcrypt.GenerateFromPassword([]byte(password), bcrypt.DefaultCost)
	return string(bytes), err
}

// Helper to check password
func checkPasswordHash(password, hash string) bool {
	return bcrypt.CompareHashAndPassword([]byte(hash), []byte(password)) == nil
}

// Helper to generate JWT
func generateJWT(userID string) (string, error) {
	secret := os.Getenv("JWT_SECRET")
	claims := jwt.MapClaims{
		"user_id": userID,
		"exp":     time.Now().Add(24 * time.Hour).Unix(),
	}
	token := jwt.NewWithClaims(jwt.SigningMethodHS256, claims)
	return token.SignedString([]byte(secret))
}

// Helper to convert database user to response model
func userToResponse(user *database.Users) database.UserResponse {
	// Handle type assertions safely
	var email string
	if user.Email != nil {
		if str, ok := user.Email.(string); ok {
			email = str
		}
	}

	var username string
	if user.Username != nil {
		if str, ok := user.Username.(string); ok {
			username = str
		}
	}

	var firstName string
	if user.First_name != nil {
		if str, ok := user.First_name.(string); ok {
			firstName = str
		}
	}

	var lastName string
	if user.Last_name != nil {
		if str, ok := user.Last_name.(string); ok {
			lastName = str
		}
	}

	return database.UserResponse{
		ID:        user.Id,
		Email:     email,
		Username:  username,
		FirstName: firstName,
		LastName:  lastName,
		CreatedAt: user.Created_at,
		UpdatedAt: user.Updated_at,
	}
}

// Users handlers
func (s *FiberServer) createUser(c *fiber.Ctx) error {
	var req database.CreateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Hash password
	hash, err := hashPassword(req.Password)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to hash password")
	}

	// Create database user
	user := database.Users{
		Email:         req.Email,
		Username:      req.Username,
		Password_hash: hash,
		First_name:    req.FirstName,
		Last_name:     req.LastName,
		Created_at:    time.Now(),
		Updated_at:    time.Now(),
	}

	// Log the user struct being created
	fmt.Printf("DEBUG: Creating user struct: %+v\n", user)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	createdUser, err := s.db.CreateUser(ctx, &user)
	if err != nil {
		fmt.Printf("DEBUG: CreateUser error: %v\n", err)
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to create user: "+err.Error())
	}

	// Invalidate users list cache
	s.cache.Del(ctx, "users:list:*")

	return c.Status(fiber.StatusCreated).JSON(fiber.Map{
		"data": userToResponse(createdUser),
	})
}

func (s *FiberServer) getUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return errorResponse(c, fiber.StatusBadRequest, "User ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to get from cache first
	cacheKey := userCacheKey(id)
	if cachedData, err := s.GetCache(ctx, cacheKey); err == nil {
		var user database.Users
		if json.Unmarshal([]byte(cachedData), &user) == nil {
			return successResponse(c, userToResponse(&user))
		}
	}

	// Get from database
	user, err := s.db.GetUserByID(ctx, id)
	if err != nil {
		return errorResponse(c, fiber.StatusNotFound, "User not found")
	}

	// Cache the user data (without password hash)
	userToCache := *user
	userToCache.Password_hash = ""
	if userData, err := json.Marshal(userToCache); err == nil {
		s.SetCache(ctx, cacheKey, string(userData), 10*time.Minute)
	}

	return successResponse(c, userToResponse(user))
}

func (s *FiberServer) listUsers(c *fiber.Ctx) error {
	limit, offset := getPaginationParams(c)

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Try to get from cache first
	cacheKey := usersListCacheKey(limit, offset)
	if cachedData, err := s.GetCache(ctx, cacheKey); err == nil {
		var users []database.Users
		if json.Unmarshal([]byte(cachedData), &users) == nil {
			// Convert to response models
			responses := make([]database.UserResponse, len(users))
			for i, user := range users {
				responses[i] = userToResponse(&user)
			}
			return successResponse(c, responses)
		}
	}

	// Get from database
	users, err := s.db.ListUsers(ctx, limit, offset)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to fetch users: "+err.Error())
	}

	// Cache the users data (without password hashes)
	usersToCache := make([]database.Users, len(users))
	for i, user := range users {
		usersToCache[i] = user
		usersToCache[i].Password_hash = ""
	}
	if usersData, err := json.Marshal(usersToCache); err == nil {
		s.SetCache(ctx, cacheKey, string(usersData), 10*time.Minute)
	}

	// Convert to response models
	responses := make([]database.UserResponse, len(users))
	for i, user := range users {
		responses[i] = userToResponse(&user)
	}

	return successResponse(c, responses)
}

func (s *FiberServer) updateUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return errorResponse(c, fiber.StatusBadRequest, "User ID is required")
	}

	var req database.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	// Get existing user
	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	existingUser, err := s.db.GetUserByID(ctx, id)
	if err != nil {
		return errorResponse(c, fiber.StatusNotFound, "User not found")
	}

	// Update fields if provided
	if req.Email != nil {
		existingUser.Email = *req.Email
	}
	if req.Username != nil {
		existingUser.Username = *req.Username
	}
	if req.FirstName != nil {
		existingUser.First_name = *req.FirstName
	}
	if req.LastName != nil {
		existingUser.Last_name = *req.LastName
	}
	existingUser.Updated_at = time.Now()

	updatedUser, err := s.db.UpdateUser(ctx, existingUser)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to update user: "+err.Error())
	}

	// Invalidate cache
	s.DeleteCache(ctx, userCacheKey(id))
	s.cache.Del(ctx, "users:list:*")

	return successResponse(c, userToResponse(updatedUser))
}

func (s *FiberServer) deleteUser(c *fiber.Ctx) error {
	id := c.Params("id")
	if id == "" {
		return errorResponse(c, fiber.StatusBadRequest, "User ID is required")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	err := s.db.DeleteUser(ctx, id)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to delete user: "+err.Error())
	}

	// Invalidate cache
	s.DeleteCache(ctx, userCacheKey(id))
	s.cache.Del(ctx, "users:list:*")

	return c.Status(fiber.StatusNoContent).Send(nil)
}

// POST /api/v1/auth/login
func (s *FiberServer) loginUser(c *fiber.Ctx) error {
	var req database.LoginRequest
	if err := c.BodyParser(&req); err != nil {
		return errorResponse(c, fiber.StatusBadRequest, "Invalid request body")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	// Find user by email
	user, err := s.db.GetUserByEmail(ctx, req.Email)
	if err != nil {
		return errorResponse(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	// Handle type assertion for password hash
	var passwordHash string
	if user.Password_hash != nil {
		if str, ok := user.Password_hash.(string); ok {
			passwordHash = str
		}
	}

	if user == nil || passwordHash == "" {
		return errorResponse(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	if !checkPasswordHash(req.Password, passwordHash) {
		return errorResponse(c, fiber.StatusUnauthorized, "Invalid credentials")
	}

	// Generate JWT
	token, err := generateJWT(user.Id)
	if err != nil {
		return errorResponse(c, fiber.StatusInternalServerError, "Failed to generate token")
	}

	response := database.LoginResponse{
		User:  userToResponse(user),
		Token: token,
	}

	return successResponse(c, response)
}
