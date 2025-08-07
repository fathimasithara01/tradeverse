package controllers

import (
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

func GetUserNotifications(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	notifs, err := notifService.GetUserNotifications(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch notifications"})
		return
	}
	c.JSON(http.StatusOK, notifs)
}

func MarkUserNotificationsRead(c *gin.Context) {
	id, _ := strconv.Atoi(c.Param("id"))
	err := notifService.MarkAsRead(uint(id))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to update notifications"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "All marked as read"})
}
