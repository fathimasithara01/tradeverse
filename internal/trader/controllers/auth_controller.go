package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/internal/trader/service"
	"github.com/gin-gonic/gin"
)

type TraderController struct {
	svc service.TraderService
}

func NewTraderController(svc service.TraderService) *TraderController {
	return &TraderController{svc: svc}
}

func (t *TraderController) Signup(c *gin.Context) {
	var req struct {
		Name        string `json:"name" binding:"required"`
		Email       string `json:"email" binding:"required,email"`
		Password    string `json:"password" binding:"required"`
		CompanyName string `json:"company_name"`
		Bio         string `json:"bio"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := t.svc.Signup(req.Name, req.Email, req.Password, req.CompanyName, req.Bio, "your_jwt_secret")
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}

func (t *TraderController) Login(c *gin.Context) {
	var req struct {
		Email    string `json:"email" binding:"required,email"`
		Password string `json:"password" binding:"required"`
	}

	if err := c.ShouldBindJSON(&req); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	token, err := t.svc.Login(req.Email, req.Password, "your_jwt_secret")
	if err != nil {
		c.JSON(http.StatusUnauthorized, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"token": token})
}
