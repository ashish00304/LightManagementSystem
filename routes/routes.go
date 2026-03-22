package routes

import (
	"github.com/gin-gonic/gin"
	"light-management/handlers"
	"light-management/repository"
	"light-management/services"
)

func SetupRouter() *gin.Engine {
	r := gin.Default()

	// Dependency injection
	repo := repository.NewMemoryLightRepository()
	service := services.NewLightService(repo)
	handler := handlers.NewLightHandler(service)

	r.GET("/lights", handler.GetLights)
	r.POST("/lights", handler.CreateLight)
	r.POST("/lights/:id/on", handler.TurnOnLight)
	r.POST("/lights/:id/off", handler.TurnOffLight)

	return r
}
