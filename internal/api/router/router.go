package router

import (
	"money-transfer/internal/api/interfaces"

	"github.com/gin-gonic/gin"

	_ "money-transfer/internal/api/docs" // swagger docs import

	swaggerFiles "github.com/swaggo/files"
	ginSwagger "github.com/swaggo/gin-swagger"
)

// NewRouter creates and configures a new router
func NewRouter(handlers []interfaces.Handler) *gin.Engine {
	router := gin.New()
	router.Use(gin.Recovery())

	// Swagger route
	router.GET("/swagger/*any", ginSwagger.WrapHandler(swaggerFiles.Handler))

	// API v1 group
	v1 := router.Group("/api/v1")

	// Register all handlers
	for _, h := range handlers {
		h.Register(v1)
	}

	return router
}
