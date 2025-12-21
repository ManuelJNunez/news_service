package news

import (
	"log/slog"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Handler struct {
	svc Service
}

func NewHandler(svc Service) *Handler {
	return &Handler{svc: svc}
}

func RegisterRoutes(rg *gin.RouterGroup, h *Handler) {
	grp := rg.Group("/news")

	grp.GET("", h.getNews)
	slog.Info("news routes registered")
}

func validateAndParseID(idStr string) (uint64, error) {
	id, err := strconv.ParseUint(idStr, 10, 64)
	if err != nil {
		return 0, err
	}

	return id, nil
}

func (h *Handler) getNews(c *gin.Context) {
	idStr := c.Query("id")
	clientIP := c.ClientIP()
	slog.Debug("article request received", slog.String("id", idStr), slog.String("client_ip", clientIP))

	// If the ID is empty, return a bad request error
	if idStr == "" {
		slog.Warn("invalid request: missing id parameter", slog.String("client_ip", clientIP))
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	// If the ID is not a valid unsigned int, return a bad request error
	id, err := validateAndParseID(idStr)
	if err != nil {
		slog.Warn("invalid request: id must be a valid number", slog.String("id", idStr), slog.String("client_ip", clientIP))
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a valid number"})
		return
	}

	// Get the article by ID from the service and handle any errors
	article, err := h.svc.GetByID(c.Request.Context(), id)

	//If there was an error getting the article, return non found error
	if err != nil {
		// Log the actual error for debugging
		slog.Error("error fetching article", slog.Uint64("id", id), slog.String("client_ip", clientIP), slog.Any("error", err))
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	// Return the article rendered as HTML
	slog.Info("article request successful", slog.Uint64("id", id), slog.String("client_ip", clientIP))
	c.HTML(http.StatusOK, "article.html", article)
}
