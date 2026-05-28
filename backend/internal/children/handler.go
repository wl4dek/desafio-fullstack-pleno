package children

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type ChildHandler struct {
	service *ChildService
}

func NewChildHandler(service *ChildService) *ChildHandler {
	return &ChildHandler{service: service}
}

func (h *ChildHandler) List(c *gin.Context) {
	filters := Filters{
		Name:         c.Query("childName"),
		Neighborhood: c.Query("neighborhood"),
		Page:         1,
		PerPage:      10,
	}

	if childName := c.Query("childName"); childName != "" {
		filters.Name = childName
	}

	if pageStr := c.Query("page"); pageStr != "" {
		if page, err := strconv.Atoi(pageStr); err == nil && page >= 1 {
			filters.Page = page
		}
	}

	if perPageStr := c.Query("per_page"); perPageStr != "" {
		if perPage, err := strconv.Atoi(perPageStr); err == nil && perPage >= 10 && perPage <= 50 {
			filters.PerPage = perPage
		}
	}

	if hasAlertStr := c.Query("has_alert"); hasAlertStr != "" {
		if val, err := strconv.ParseBool(hasAlertStr); err == nil {
			filters.HasAlert = &val
		}
	}

	if reviewedStr := c.Query("reviewed"); reviewedStr != "" {
		if val, err := strconv.ParseBool(reviewedStr); err == nil {
			filters.Reviewed = &val
		}
	}

	result, err := h.service.List(c.Request.Context(), filters)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, result)
}

func (h *ChildHandler) GetByID(c *gin.Context) {
	id := c.Param("id")
	child, err := h.service.GetByID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if child == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "child not found"})
		return
	}

	c.JSON(http.StatusOK, child)
}

func (h *ChildHandler) MarkReviewed(c *gin.Context) {
	id := c.Param("id")
	reviewedBy, _ := c.Get("preferred_username")
	reviewedByStr, _ := reviewedBy.(string)

	err := h.service.MarkReviewed(c.Request.Context(), id, reviewedByStr)
	if err != nil {
		if errors.Is(err, ErrChildNotFound) {
			c.JSON(http.StatusNotFound, gin.H{"error": "child not found"})
			return
		}
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "review registered"})
}

func (h *ChildHandler) GetAreasByChildID(c *gin.Context) {
	id := c.Param("id")
	areas, err := h.service.GetAreasByChildID(c.Request.Context(), id)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	if areas == nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "child not found"})
		return
	}

	c.JSON(http.StatusOK, areas)
}

func (h *ChildHandler) ListNeighborhood(c *gin.Context) {
	neighborhoods, err := h.service.ListNeighborhood(c.Request.Context())
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "internal error"})
		return
	}

	c.JSON(http.StatusOK, neighborhoods)
}
