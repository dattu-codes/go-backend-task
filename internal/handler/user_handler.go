package handler

import (
	"errors"
	"strconv"
	"time"

	"go-backend-task/internal/logger"
	"go-backend-task/internal/models"
	"go-backend-task/internal/repository"
	"go-backend-task/internal/service"

	"github.com/go-playground/validator/v10"
	"github.com/gofiber/fiber/v2"
	"github.com/jackc/pgx/v5"
	"go.uber.org/zap"
)

// UserHandler hosts the HTTP controller methods for User REST endpoints.
type UserHandler struct {
	repo     repository.UserRepository
	validate *validator.Validate
}

// NewUserHandler instantiates a UserHandler injects repository dependency and validator.
func NewUserHandler(repo repository.UserRepository) *UserHandler {
	return &UserHandler{
		repo:     repo,
		validate: validator.New(),
	}
}

// Create inserts a user into the DB after parsing and validating the payload.
// Returns HTTP 201 Created on success.
func (h *UserHandler) Create(c *fiber.Ctx) error {
	reqID, _ := c.Locals("requestId").(string)
	var req models.CreateUserRequest

	if err := c.BodyParser(&req); err != nil {
		logger.Log.Warn("JSON parsing failed", zap.String("request_id", reqID), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if err := h.validate.Struct(req); err != nil {
		logger.Log.Warn("Validation failed", zap.String("request_id", reqID), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	// datetime=2006-01-02 validator tag guarantees successful parsing here
	dobTime, _ := time.Parse("2006-01-02", req.DOB)

	user, err := h.repo.CreateUser(c.UserContext(), req.Name, dobTime)
	if err != nil {
		logger.Log.Error("Database insertion failed", zap.String("request_id", reqID), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to create user"})
	}

	age := service.CalculateAge(user.Dob.Time, time.Now())

	return c.Status(fiber.StatusCreated).JSON(models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Time.Format("2006-01-02"),
		Age:  age,
	})
}

// Get queries a single user by ID. Returns HTTP 200 OK or HTTP 404 Not Found.
func (h *UserHandler) Get(c *fiber.Ctx) error {
	reqID, _ := c.Locals("requestId").(string)
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID must be a positive integer"})
	}

	user, err := h.repo.GetUserByID(c.UserContext(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		logger.Log.Error("Database fetch failed", zap.String("request_id", reqID), zap.Int("id", id), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to fetch user"})
	}

	age := service.CalculateAge(user.Dob.Time, time.Now())

	return c.Status(fiber.StatusOK).JSON(models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Time.Format("2006-01-02"),
		Age:  age,
	})
}

// Update modifies the details of an existing user. Returns HTTP 200 OK or HTTP 404 Not Found.
func (h *UserHandler) Update(c *fiber.Ctx) error {
	reqID, _ := c.Locals("requestId").(string)
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID must be a positive integer"})
	}

	var req models.UpdateUserRequest
	if err := c.BodyParser(&req); err != nil {
		logger.Log.Warn("JSON parsing failed on update", zap.String("request_id", reqID), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "Invalid request payload"})
	}

	if err := h.validate.Struct(req); err != nil {
		logger.Log.Warn("Validation failed on update", zap.String("request_id", reqID), zap.Error(err))
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": err.Error()})
	}

	dobTime, _ := time.Parse("2006-01-02", req.DOB)

	user, err := h.repo.UpdateUser(c.UserContext(), int32(id), req.Name, dobTime)
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		logger.Log.Error("Database update failed", zap.String("request_id", reqID), zap.Int("id", id), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to update user"})
	}

	age := service.CalculateAge(user.Dob.Time, time.Now())

	return c.Status(fiber.StatusOK).JSON(models.UserResponse{
		ID:   user.ID,
		Name: user.Name,
		DOB:  user.Dob.Time.Format("2006-01-02"),
		Age:  age,
	})
}

// Delete drops a user from the DB. Returns HTTP 204 No Content or HTTP 404 Not Found.
func (h *UserHandler) Delete(c *fiber.Ctx) error {
	reqID, _ := c.Locals("requestId").(string)
	idStr := c.Params("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id <= 0 {
		return c.Status(fiber.StatusBadRequest).JSON(fiber.Map{"error": "User ID must be a positive integer"})
	}

	// Verifies first that the user exists (crucial for returning explicit 404 errors)
	_, err = h.repo.GetUserByID(c.UserContext(), int32(id))
	if err != nil {
		if errors.Is(err, pgx.ErrNoRows) {
			return c.Status(fiber.StatusNotFound).JSON(fiber.Map{"error": "User not found"})
		}
		logger.Log.Error("User verification failed for deletion", zap.String("request_id", reqID), zap.Int("id", id), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to verify user"})
	}

	err = h.repo.DeleteUser(c.UserContext(), int32(id))
	if err != nil {
		logger.Log.Error("Database delete failed", zap.String("request_id", reqID), zap.Int("id", id), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to delete user"})
	}

	return c.SendStatus(fiber.StatusNoContent)
}

// List fetches multiple users with page/limit parameters. Returns HTTP 200 OK.
func (h *UserHandler) List(c *fiber.Ctx) error {
	reqID, _ := c.Locals("requestId").(string)

	// Fetch query parameters for pagination, with defaults
	pageStr := c.Query("page", "1")
	limitStr := c.Query("limit", "10")

	page, err := strconv.Atoi(pageStr)
	if err != nil || page < 1 {
		page = 1
	}

	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 || limit > 100 {
		limit = 10
	}

	offset := (page - 1) * limit

	users, err := h.repo.ListUsers(c.UserContext(), int32(limit), int32(offset))
	if err != nil {
		logger.Log.Error("Database list query failed", zap.String("request_id", reqID), zap.Error(err))
		return c.Status(fiber.StatusInternalServerError).JSON(fiber.Map{"error": "Failed to list users"})
	}

	now := time.Now()
	responses := make([]models.UserResponse, len(users))
	for i, user := range users {
		responses[i] = models.UserResponse{
			ID:   user.ID,
			Name: user.Name,
			DOB:  user.Dob.Time.Format("2006-01-02"),
			Age:  service.CalculateAge(user.Dob.Time, now),
		}
	}

	return c.Status(fiber.StatusOK).JSON(responses)
}
