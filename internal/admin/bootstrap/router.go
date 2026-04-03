package bootstrap

import (
	"github.com/fathimasithara01/tradeverse/config"
	"github.com/fathimasithara01/tradeverse/internal/admin/router"
	"github.com/gin-gonic/gin"
	"gorm.io/gorm"
)

func InitRouter(s *Services, cfg *config.Config, db *gorm.DB) *gin.Engine {
	r := gin.Default()

	ctrls := InitControllers(s)

	router.WireAdminRoutes(
		r,
		cfg,
		ctrls.Auth,
		ctrls.Dashboard,
		ctrls.User,
		ctrls.Role,
		ctrls.Permission,
		ctrls.Activity,
		s.Role,
		ctrls.AdminWallet,
		ctrls.Subscription,
		ctrls.Transaction,
		db,
		ctrls.Signal,
		ctrls.Commission,
		ctrls.WebConfiguration,
	)

	return r
}
