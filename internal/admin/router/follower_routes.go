package router

import (
	"github.com/fathimasithara01/tradeverse/internal/admin/controllers"
	"github.com/fathimasithara01/tradeverse/internal/admin/middleware"
	"github.com/fathimasithara01/tradeverse/pkg/config"
	"github.com/gin-gonic/gin"
)

func WireFollowerRoutes(r *gin.Engine, copyCtrl *controllers.CopyController, cfg *config.Config) {
	protected := r.Group("/api/copy")
	protected.Use(middleware.JWTMiddleware(cfg))
	{
		protected.GET("/status/:masterID", copyCtrl.GetCopyStatus)
		protected.POST("/start", copyCtrl.StartCopying)
		protected.POST("/stop", copyCtrl.StopCopying)
	}
}
