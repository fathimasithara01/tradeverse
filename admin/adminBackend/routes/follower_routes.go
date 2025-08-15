package routes

import (
	"github.com/fathimasithara01/tradeverse/controllers"
	"github.com/fathimasithara01/tradeverse/middleware"
	"github.com/gin-gonic/gin"
)

func WireFollowerRoutes(r *gin.Engine, copyCtrl *controllers.CopyController) {
	protected := r.Group("/api/copy")
	protected.Use(middleware.JWTMiddleware())
	{
		protected.GET("/status/:masterID", copyCtrl.GetCopyStatus)
		protected.POST("/start", copyCtrl.StartCopying)
		protected.POST("/stop", copyCtrl.StopCopying)
	}
}
