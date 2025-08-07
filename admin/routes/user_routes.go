package routes

import (
	"github.com/fathimasithara01/tradeverse/admin/controllers"
	"github.com/fathimasithara01/tradeverse/admin/middleware"
	"github.com/gin-gonic/gin"
)

func UserRoutes(r *gin.Engine) {
	user := r.Group("/user")
	{
		user.GET("/users/:id/notifications", middleware.JWTMiddleware(), controllers.GetUserNotifications)
		user.PUT("/users/:id/notifications/read-all", middleware.JWTMiddleware(), controllers.MarkUserNotificationsRead)

		user.POST("/payments/checkout", controllers.CreateCheckout)
		user.POST("/payments/verify", controllers.VerifyPayment)

	}

}
