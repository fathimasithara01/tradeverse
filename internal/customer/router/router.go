package router

import (
	"github.com/fathimasithara01/tradeverse/internal/customer/controllers"
	"github.com/fathimasithara01/tradeverse/internal/customer/middleware"
	"github.com/fathimasithara01/tradeverse/pkg/config"
	"github.com/gin-gonic/gin"
)

func SetupRouter(
	cfg *config.Config,
	authController *controllers.AuthController,
	profileController *controllers.ProfileController,
) *gin.Engine {
	r := gin.Default()

	public := r.Group("/api/v1")
	{
		public.POST("/signup", authController.Signup)
		public.POST("/login", authController.Login)
	}

	protected := r.Group("/api/v1")
	protected.Use(middleware.AuthMiddleware(cfg.JWTSecret))
	{
		protected.GET("/profile", profileController.GetProfile)
		protected.PUT("/profile", profileController.UpdateProfile)
		protected.DELETE("/account", profileController.DeleteAccount)
	}

	return r
}
