package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/service"
	"github.com/gin-gonic/gin"
)

type SignalController struct {
	LiveSignalSvc *service.LiveSignalService
}

func NewSignalController(liveSignalSvc *service.LiveSignalService) *SignalController {
	return &SignalController{LiveSignalSvc: liveSignalSvc}
}

func (ctrl *SignalController) ShowLiveSignalsPage(c *gin.Context) {
	c.HTML(http.StatusOK, "live_signals.html", nil)
}

// GetLiveSignals provides the JSON data for the frontend JavaScript.
func (ctrl *SignalController) GetLiveSignals(c *gin.Context) {
	signals, err := ctrl.LiveSignalSvc.GetLiveSignals()
	if err != nil {
		// Send the error response.
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve live signals"})
		// CRITICAL FIX: Immediately stop the function execution.
		return
	}

	// If the service returns `nil` (no signals), send an empty array `[]`
	// which is safer for JavaScript to handle than `null`.
	if signals == nil {
		signals = make([]service.TraderSignal, 0)
	}

	// This line will now only run if there was NO error.
	c.JSON(http.StatusOK, signals)
}
