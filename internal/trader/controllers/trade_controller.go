package controllers

import (
	"errors"
	"net/http"
	"strconv"

	"github.com/fathimasithara01/tradeverse/internal/trader/service"
	"github.com/fathimasithara01/tradeverse/pkg/models"
	"github.com/fathimasithara01/tradeverse/pkg/utils/constants"
	"github.com/fathimasithara01/tradeverse/pkg/utils/response"
	"github.com/gin-gonic/gin"
)

type TradeController struct {
	tradeService service.TradeService
}

func NewTradeController(tradeService service.TradeService) *TradeController {
	return &TradeController{
		tradeService: tradeService,
	}
}

func getTraderIDAndCheckRole(c *gin.Context) (uint, error) {
	userIDVal, exists := c.Get("userID")
	if !exists {
		return 0, errors.New("userID not found in context")
	}
	traderID, ok := userIDVal.(uint)
	if !ok {
		return 0, errors.New("userID in context is not of type uint")
	}

	userRoleVal, exists := c.Get("userRole")
	if !exists {
		return 0, errors.New("userRole not found in context")
	}
	userRole, ok := userRoleVal.(string)
	if !ok {
		return 0, errors.New("userRole in context is not of type string")
	}

	if userRole != string(models.RoleTrader) {
		return 0, errors.New("forbidden: only traders can access this resource")
	}

	return traderID, nil
}

func (ctrl *TradeController) CreateTrade(c *gin.Context) {
	traderID, err := getTraderIDAndCheckRole(c)
	if err != nil {
		if err.Error() == "forbidden: only traders can access this resource" {
			response.Error(c, http.StatusForbidden, err.Error())
		} else {
			response.Error(c, http.StatusUnauthorized, "Unauthorized") 
		}
		return
	}

	var input models.TradeInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	trade, err := ctrl.tradeService.CreateTrade(c.Request.Context(), traderID, input)
	if err != nil {
		if errors.Is(err, constants.ErrForbidden) {
			response.Error(c, http.StatusForbidden, err.Error())
		} else {
			response.Error(c, http.StatusInternalServerError, err.Error())
		}
		return
	}

	response.Success(c, http.StatusCreated, "Trade created successfully", trade)
}

func (ctrl *TradeController) GetTraderTrades(c *gin.Context) {
	traderID, err := getTraderIDAndCheckRole(c)
	if err != nil {
		if err.Error() == "forbidden: only traders can access this resource" {
			response.Error(c, http.StatusForbidden, err.Error())
		} else {
			response.Error(c, http.StatusUnauthorized, "Unauthorized")
		}
		return
	}

	page, _ := strconv.Atoi(c.DefaultQuery("page", "1"))
	limit, _ := strconv.Atoi(c.DefaultQuery("limit", "10"))

	tradeListResponse, err := ctrl.tradeService.GetTraderTrades(c.Request.Context(), traderID, page, limit)
	if err != nil {
		response.Error(c, http.StatusInternalServerError, err.Error())
		return
	}

	response.Success(c, http.StatusOK, "Trader trade history retrieved successfully", tradeListResponse)
}

func (ctrl *TradeController) UpdateTradeStatus(c *gin.Context) {
	traderID, err := getTraderIDAndCheckRole(c)
	if err != nil {
		if err.Error() == "forbidden: only traders can access this resource" {
			response.Error(c, http.StatusForbidden, err.Error())
		} else {
			response.Error(c, http.StatusUnauthorized, "Unauthorized")
		}
		return
	}

	tradeIDStr := c.Param("id")
	tradeID, err := strconv.ParseUint(tradeIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid trade ID")
		return
	}

	var input models.TradeUpdateInput
	if err := c.ShouldBindJSON(&input); err != nil {
		response.Error(c, http.StatusBadRequest, err.Error())
		return
	}

	updatedTrade, err := ctrl.tradeService.UpdateTradeStatus(c.Request.Context(), traderID, uint(tradeID), input)
	if err != nil {
		if errors.Is(err, constants.ErrNotFound("Trade")) {
			response.Error(c, http.StatusNotFound, err.Error())
		} else if errors.Is(err, constants.ErrForbidden) || err.Error() == "cannot modify a trade that is not open or pending" { // Catch business logic errors as Bad Request
			response.Error(c, http.StatusForbidden, err.Error())
		} else {
			response.Error(c, http.StatusBadRequest, err.Error()) // Use 400 for general business logic errors
		}
		return
	}

	response.Success(c, http.StatusOK, "Trade updated successfully", updatedTrade)
}

func (ctrl *TradeController) RemoveTrade(c *gin.Context) {
	traderID, err := getTraderIDAndCheckRole(c)
	if err != nil {
		if err.Error() == "forbidden: only traders can access this resource" {
			response.Error(c, http.StatusForbidden, err.Error())
		} else {
			response.Error(c, http.StatusUnauthorized, "Unauthorized")
		}
		return
	}

	tradeIDStr := c.Param("id")
	tradeID, err := strconv.ParseUint(tradeIDStr, 10, 32)
	if err != nil {
		response.Error(c, http.StatusBadRequest, "Invalid trade ID")
		return
	}

	err = ctrl.tradeService.RemoveTrade(c.Request.Context(), traderID, uint(tradeID))
	if err != nil {
		if errors.Is(err, constants.ErrNotFound("Trade")) {
			response.Error(c, http.StatusNotFound, err.Error())
		} else if errors.Is(err, constants.ErrForbidden) || err.Error() == "only pending trades can be removed" { // Catch business logic errors as Bad Request
			response.Error(c, http.StatusForbidden, err.Error())
		} else {
			response.Error(c, http.StatusBadRequest, err.Error()) // Use 400 for general business logic errors
		}
		return
	}

	c.Status(http.StatusNoContent)
}
