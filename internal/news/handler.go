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
	slog.Debug("news request received", slog.String("id", idStr), slog.String("client_ip", clientIP))

	if idStr == "" {
		slog.Warn("invalid request: missing id parameter", slog.String("client_ip", clientIP))
		c.JSON(http.StatusBadRequest, gin.H{"error": "id parameter is required"})
		return
	}

	id, err := validateAndParseID(idStr)
	if err != nil {
		slog.Warn("invalid request: id must be a valid number", slog.String("id", idStr), slog.String("client_ip", clientIP))
		c.JSON(http.StatusBadRequest, gin.H{"error": "id must be a valid number"})
		return
	}

	article, err := h.svc.GetByID(c.Request.Context(), id)

	if err != nil {
		if err == ErrNewsNotFound {
			slog.Warn("news not found", slog.Uint64("id", id), slog.String("client_ip", clientIP))
			c.JSON(http.StatusNotFound, gin.H{"error": "news not found"})
			return
		}
		slog.Error("error fetching news", slog.Uint64("id", id), slog.String("client_ip", clientIP), slog.Any("error", err))
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	slog.Info("news request successful", slog.Uint64("id", id), slog.String("client_ip", clientIP))
	c.HTML(http.StatusOK, "article.html", article)
}
