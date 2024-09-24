package controllers

import (
	"GeekReward/inernal/app/services"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

type GeekController struct {
	geekService services.GeekService
}

func NewGeekController(geekService services.GeekService) *GeekController {
	return &GeekController{geekService: geekService}
}

// GetTopGeeks 获取排名前的极客用户
func (ctl *GeekController) GetTopGeeks(c *gin.Context) {
	limitStr := c.DefaultQuery("limit", "10")
	limit, err := strconv.Atoi(limitStr)
	if err != nil || limit < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid limit"})
		return
	}

	geeks, err := ctl.geekService.GetTopGeeks(limit)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch geeks"})
		return
	}

	c.JSON(http.StatusOK, geeks)
}

// GetGeekByID 获取指定ID的极客用户
func (ctl *GeekController) GetGeekByID(c *gin.Context) {
	idStr := c.Param("id")
	id, err := strconv.Atoi(idStr)
	if err != nil || id < 1 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid ID"})
		return
	}

	geek, err := ctl.geekService.GetGeekByID(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch geek"})
		return
	}

	c.JSON(http.StatusOK, geek)
}
