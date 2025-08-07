package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/admin/models"
	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/service"
	"github.com/gin-gonic/gin"
)

var announcementService = service.AnnouncementService{
	Repo: repository.AnnouncementRepository{},
}

func CreateAnnouncement(c *gin.Context) {
	var a models.Announcement
	if err := c.ShouldBindJSON(&a); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}
	saved, err := announcementService.Create(a)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Could not create announcement"})
		return
	}
	c.JSON(http.StatusCreated, saved)
}

func GetAllAnnouncements(c *gin.Context) {
	all, err := announcementService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch announcements"})
		return
	}
	c.JSON(http.StatusOK, all)
}
