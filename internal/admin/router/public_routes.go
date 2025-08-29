package routes

import (
	"github.com/fathimasithara01/tradeverse/internal/admin/controllers"
	"github.com/gin-gonic/gin"
)

func WirePublicRoutes(r *gin.Engine, authCtrl *controllers.AuthController, signalCtrl *controllers.SignalController) {
	r.GET("/login", authCtrl.ShowLoginPage)
	r.POST("/login", authCtrl.LoginUser)

	r.GET("/signals", signalCtrl.ShowLiveSignalsPage)
	r.GET("/api/signals", signalCtrl.GetLiveSignals)

	r.GET("/register/customer", authCtrl.ShowCustomerRegisterPage)
	r.POST("/register/customer", authCtrl.RegisterCustomer)
	r.GET("/register/trader", authCtrl.ShowTraderRegisterPage)
	r.POST("/register/trader", authCtrl.RegisterTrader)
}
