package summary

import (
	"net/http"

	"github.com/gin-gonic/gin"
)

type SummaryHandler struct {
	service *SummaryService
}

func NewSummaryHandler(service *SummaryService) *SummaryHandler {
	return &SummaryHandler{service: service}
}

func (h *SummaryHandler) GetSummary(c *gin.Context) {
	summary, err := h.service.GetSummary(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, summary)
}
