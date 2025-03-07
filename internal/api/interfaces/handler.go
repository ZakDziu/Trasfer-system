package interfaces

import "github.com/gin-gonic/gin"

// Handler represents a common interface for all handlers
type Handler interface {
	Register(group *gin.RouterGroup)
}
