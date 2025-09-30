package controllers

// import (
// 	"net/http"

// 	"github.com/fathimasithara01/tradeverse/internal/admin/service"
// 	"github.com/gin-gonic/gin"
// )

// type SignalController struct {
// 	LiveSignalSvc service.ILiveSignalService
// }

// func NewSignalController(liveSignalSvc service.ILiveSignalService) *SignalController {
// 	return &SignalController{LiveSignalSvc: liveSignalSvc}
// }

// func (ctrl *SignalController) ShowLiveSignalsPage(c *gin.Context) {
// 	c.HTML(http.StatusOK, "live_signals.html", nil)
// }

// func (ctrl *SignalController) GetLiveSignals(c *gin.Context) {
// 	signals, err := ctrl.LiveSignalSvc.GetLiveSignals()
// 	if err != nil {
// 		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve live signals"})
// 		return
// 	}

// 	if signals == nil {
// 		signals = make([]service.TraderSignal, 0)
// 	}

// 	c.JSON(http.StatusOK, signals)
// }
