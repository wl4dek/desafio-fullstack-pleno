package statistics

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type StatisticsHandler struct {
	service *StatisticsService
}

func NewStatisticsHandler(service *StatisticsService) *StatisticsHandler {
	return &StatisticsHandler{service: service}
}

func (h *StatisticsHandler) GetStatistics(c *gin.Context) {
	result, err := h.service.GetStatistics(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *StatisticsHandler) GetSummary(c *gin.Context) {
	summary, err := h.service.GetSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, summary)
}
