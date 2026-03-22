package handlers

import (
	"fmt"
	"net/http"

	"github.com/gin-gonic/gin"
	"light-management/models"
	"light-management/services"
)

type LightHandler struct {
	service services.LightService
}

func NewLightHandler(service services.LightService) *LightHandler {
	return &LightHandler{service: service}
}

func (h *LightHandler) GetLights(c *gin.Context) {
	lights := h.service.GetAll()
	// Return empty array instead of null if no lights exist
	if lights == nil {
		lights = []models.Light{}
	}
	c.JSON(http.StatusOK, lights)
}

func (h *LightHandler) CreateLight(c *gin.Context) {
	var input struct {
		Name string `json:"name" binding:"required"`
	}
	if err := c.ShouldBindJSON(&input); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Name is required"})
		return
	}

	light := h.service.Create(input.Name)
	c.JSON(http.StatusCreated, light)
}

func (h *LightHandler) TurnOnLight(c *gin.Context) {
	id := c.Param("id")

	var req models.TurnOnRequest
	if err := c.ShouldBindJSON(&req); err != nil && err.Error() != "EOF" && c.Request.ContentLength > 0 {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	light, err := h.service.TurnOn(id, req.Duration)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Light %s turned ON", id), "light": light})
}

func (h *LightHandler) TurnOffLight(c *gin.Context) {
	id := c.Param("id")

	light, err := h.service.TurnOff(id)
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": fmt.Sprintf("Light %s turned OFF", id), "light": light})
}
