package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/service"
	"github.com/gin-gonic/gin"
)

var subscriptionService = service.SubscriptionService{
	Repo: repository.SubscriptionRepository{},
}

func GetAllSubscriptions(c *gin.Context) {
	subs, err := subscriptionService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch subscriptions"})
		return
	}
	c.JSON(http.StatusOK, subs)
}
