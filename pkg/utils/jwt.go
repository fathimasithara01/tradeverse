package utils

import (
	"errors"

	"github.com/gin-gonic/gin"
)

const UserIDKey = "userID"

func GetUserIDFromContext(c *gin.Context) (uint, error) {
	userID, exists := c.Get(UserIDKey)
	if !exists {
		return 0, errors.New("user ID not found in context")
	}

	id, ok := userID.(uint)
	if !ok {
		return 0, errors.New("invalid user ID type in context")
	}
	return id, nil
}

func SetUserIDInContext(c *gin.Context, userID uint) {
	c.Set(UserIDKey, userID)
}
