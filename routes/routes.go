package routes

import (
	"github.com/draco121/authorizationservice/controllers"
	"github.com/draco121/common/utils"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(controllers controllers.Controllers, router *gin.Engine) {
	utils.Logger.Info("Registering routes...")
	v1 := router.Group("/v1")
	v1.POST("/authorize", controllers.Authorize)
	utils.Logger.Info("Registered routes...")
}
