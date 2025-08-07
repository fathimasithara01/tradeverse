package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/service"
	"github.com/gin-gonic/gin"
)

var notifService = service.NotificationService{
	Repo: repository.NotificationRepository{},
}

func SendNotification(c *gin.Context) {
	var body struct {
		UserID  uint   `json:"user_id"`
		Title   string `json:"title"`
		Message string `json:"message"`
		Type    string `json:"type"` // info, alert, warning
	}
	if err := c.ShouldBindJSON(&body); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request"})
		return
	}
	err := notifService.Send(body.UserID, body.Title, body.Message, body.Type)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to send notification"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "Notification sent"})
}
