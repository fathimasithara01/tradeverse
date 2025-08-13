package routes

import (
	"github.com/fathimasithara01/tradeverse/controllers"
	"github.com/gin-gonic/gin"
)

func WirePublicRoutes(r *gin.Engine, authCtrl *controllers.AuthController) {
	// Universal Login for all user types
	r.GET("/login", authCtrl.ShowLoginPage)
	r.POST("/login", authCtrl.LoginUser)

	// Specialized Registration for each public role
	r.GET("/register/customer", authCtrl.ShowCustomerRegisterPage)
	r.POST("/register/customer", authCtrl.RegisterCustomer)
	r.GET("/register/trader", authCtrl.ShowTraderRegisterPage)
	r.POST("/register/trader", authCtrl.RegisterTrader)
}
