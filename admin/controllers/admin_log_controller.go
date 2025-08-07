package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/service"
	"github.com/gin-gonic/gin"
)

var logService = service.LogService{
	Repo: repository.LogRepository{},
}

func GetAllLogs(c *gin.Context) {
	logs, err := logService.GetAll()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve logs"})
		return
	}
	c.JSON(http.StatusOK, logs)
}
