package route

import (
	"github.com/gin-gonic/gin"
	"user-management/internal/api/handler"
)

func SetupUserRoutes(router *gin.RouterGroup, userHandler *handler.UserHandler) {
	users := router.Group("/users")

	users.GET("/:id", userHandler.GetProfile)
}
