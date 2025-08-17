package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/models"
	"github.com/fathimasithara01/tradeverse/service"
	"github.com/gin-gonic/gin"
)

type ActivityController struct{ ActivitySvc *service.ActivityService }

func NewActivityController(activitySvc *service.ActivityService) *ActivityController {
	return &ActivityController{ActivitySvc: activitySvc}
}

func (ctrl *ActivityController) ShowLiveCopyingPage(c *gin.Context) {
	c.HTML(http.StatusOK, "live_copying.html", nil)
}
func (ctrl *ActivityController) ShowTradeErrorsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "trade_errors.html", nil)
}

func (ctrl *ActivityController) GetActiveSessions(c *gin.Context) {
	sessions, err := ctrl.ActivitySvc.GetActiveSessions()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get live sessions"})
		return
	}
	if sessions == nil {
		sessions = make([]models.CopySession, 0)
	}
	c.JSON(http.StatusOK, sessions)
}

func (ctrl *ActivityController) GetTradeLogs(c *gin.Context) {
	logs, err := ctrl.ActivitySvc.GetRecentTradeLogs()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to get trade logs"})
		return
	}
	if logs == nil {
		logs = make([]models.TradeLog, 0)
	}
	c.JSON(http.StatusOK, logs)
}
