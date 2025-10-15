package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/trader/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/gin-gonic/gin"
)

type SignalController struct {
	signalService service.ISignalService
}

func NewSignalController(signalService service.ISignalService) *SignalController {
	return &SignalController{signalService: signalService}
}

func (ctrl *SignalController) CreateSignal(c *gin.Context) {
	var req models.Signal
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	signal, err := ctrl.signalService.CreateSignal(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create signal"})
		return
	}
	

	c.JSON(http.StatusCreated, signal)
}

func (ctrl *SignalController) GetAllSignals(c *gin.Context) {
	signals, err := ctrl.signalService.GetAllSignals(c)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch signals"})
		return
	}

	c.JSON(http.StatusOK, signals)
}

func (ctrl *SignalController) GetSignalByID(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signal ID"})
		return
	}

	signal, err := ctrl.signalService.GetSignalByID(c, uint(id))
	if err != nil {
		c.JSON(http.StatusNotFound, gin.H{"error": "signal not found"})
		return
	}

	c.JSON(http.StatusOK, signal)
}

// PUT /trader/signals/:id
func (ctrl *SignalController) UpdateSignal(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signal ID"})
		return
	}

	var req models.Signal
	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	req.ID = uint(id)

	signal, err := ctrl.signalService.UpdateSignal(c, &req)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to update signal"})
		return
	}

	c.JSON(http.StatusOK, signal)
}

// DELETE /trader/signals/:id
func (ctrl *SignalController) DeleteSignal(c *gin.Context) {
	idParam := c.Param("id")
	id, err := strconv.ParseUint(idParam, 10, 32)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "invalid signal ID"})
		return
	}

	if err := ctrl.signalService.DeleteSignal(c, uint(id)); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "failed to delete signal"})
		return
	}

	c.JSON(http.StatusOK, gin.H{"message": "signal deleted successfully"})
}
