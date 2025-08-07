package controllers

import (
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/admin/repository"
	"github.com/fathimasithara01/tradeverse/admin/service"
	"github.com/gin-gonic/gin"
)

var userService = service.UserService{
	Repo: repository.UserRepository{},
}

func GetAllUsers(c *gin.Context) {
	users, err := userService.GetAllUsers()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to fetch users"})
		return
	}
	c.JSON(http.StatusOK, users)
}

func BanUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, _ := strconv.Atoi(idParam)

	err := userService.BanUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to ban user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User banned successfully"})
}

func UnbanUser(c *gin.Context) {
	idParam := c.Param("id")
	userID, _ := strconv.Atoi(idParam)

	err := userService.UnbanUser(uint(userID))
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Unable to unban user"})
		return
	}
	c.JSON(http.StatusOK, gin.H{"message": "User unbanned successfully"})
}
