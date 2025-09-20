package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/trader/service"
	"github.com/gin-gonic/gin"
)

type SubscriberController struct {
	svc service.SubscriberService
}

func NewSubscriberController(svc service.SubscriberService) *SubscriberController {
	return &SubscriberController{svc: svc}
}

func (c *SubscriberController) ListSubscribers(ctx *gin.Context) {
	traderIDStr := ctx.GetString("userID") // from JWT
	traderID, _ := strconv.ParseUint(traderIDStr, 10, 64)

	subs, err := c.svc.ListSubscribers(uint(traderID))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	ctx.JSON(http.StatusOK, subs)
}

func (c *SubscriberController) GetSubscriber(ctx *gin.Context) {
	idParam := ctx.Param("id")
	id, _ := strconv.ParseUint(idParam, 10, 64)

	sub, err := c.svc.GetSubscriber(uint(id))
	if err != nil {
		ctx.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	if sub == nil {
		ctx.JSON(http.StatusNotFound, gin.H{"error": "subscriber not found"})
		return
	}

	ctx.JSON(http.StatusOK, sub)
}
