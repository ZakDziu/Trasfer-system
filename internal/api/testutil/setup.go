package testutil

import (
	"money-transfer/internal/api/interfaces"
	"money-transfer/internal/api/router"

	"github.com/gin-gonic/gin"
)

// SetupTestRouter creates a router for testing
func SetupTestRouter(handlers []interfaces.Handler) *gin.Engine {
	gin.SetMode(gin.TestMode)
	return router.NewRouter(handlers)
}
