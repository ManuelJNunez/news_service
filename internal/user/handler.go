package user

import (
	"log/slog"
	"net/http"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	grp := rg.Group("")

	grp.GET("/login", h.LoginGet)
	grp.POST("/login", h.LoginPost)
	grp.POST("/user/register", h.Register)
	slog.Info("user routes registered")
}

func (h *Handler) LoginGet(c *gin.Context) {
	c.HTML(http.StatusOK, "login.html", nil)
}

func (h *Handler) LoginPost(c *gin.Context) {
	// Get credentials from request body
	var credentials map[string]any
	if err := c.ShouldBindJSON(&credentials); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}

	// Extract only the relevant fields
	loginData := map[string]any{
		"username": credentials["username"],
		"password": credentials["password"],
	}

	// Find user with provided credentials
	user, err := h.svc.FindOne(c.Request.Context(), loginData)
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "Invalid credentials"})
		return
	}

	// Send successful response
	slog.Info("login successful", slog.String("username", user.Username))
	c.JSON(http.StatusOK, gin.H{
		"message":  "Login successful",
		"username": user.Username,
	})
}

func (h *Handler) Register(c *gin.Context) {
	// Get credentials from request body
	var input LoginInput
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	slog.Debug("register attempt", slog.String("username", input.Username))

	// Check username and password are not empty
	if input.Username == "" || input.Password == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Username and password are required"})
		return
	}

	// Create user
	user, err := h.svc.Create(c.Request.Context(), input)
	if err != nil {
		if err == ErrUserAlreadyExists {
			slog.Warn("user already exists", slog.String("username", input.Username))
			c.JSON(http.StatusConflict, gin.H{"error": "User already exists"})
			return
		}
		slog.Error("register error", slog.String("username", input.Username), slog.Any("error", err))
		c.JSON(http.StatusBadRequest, gin.H{"error": "Failed to create user"})
		return
	}

	// Send successful response
	slog.Info("user registered successfully", slog.String("username", user.Username))
	c.JSON(http.StatusCreated, gin.H{
		"message": "User registered successfully",
		"user":    user,
	})
}
