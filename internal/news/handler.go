package news

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
	grp := rg.Group("/news")

	grp.GET("", h.getNews)
	slog.Info("news routes registered")
}

func (h *Handler) getNews(c *gin.Context) {
	id := c.Query("id")
	clientIP := c.ClientIP()
	slog.Debug("article request received", slog.String("id", id), slog.String("client_ip", clientIP))

	// If the ID is empty, return a bad request error
	if id == "" {
		slog.Warn("invalid request: missing id parameter", slog.String("client_ip", clientIP))
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	// Get the article by ID from the service and handle any errors
	article, err := h.svc.GetByID(c.Request.Context(), id)

	//If there was an error getting the article, return not found error
	if err != nil {
		slog.Error("error fetching article", slog.String("id", id), slog.String("client_ip", clientIP), slog.Any("error", err))
		c.JSON(http.StatusNotFound, gin.H{"error": "article not found"})
		return
	}

	// Return the article rendered as HTML
	slog.Info("article request successful", slog.String("id", id), slog.String("client_ip", clientIP))
	c.HTML(http.StatusOK, "article.html", article)
}
