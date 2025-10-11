package controllers

// import (
// 	"errors"
// 	"net/http"
// 	"strconv"

// 	"github.com/fathimasithara01/tradeverse/internal/customer/service"
// 	"github.com/gin-gonic/gin"
// 	"gorm.io/gorm"
// )

// type CustomerSignalController struct {
// 	CustomerSignalSvc service.ICustomerSignalService
// }

// func NewCustomerSignalController(svc service.ICustomerSignalService) *CustomerSignalController {
// 	return &CustomerSignalController{CustomerSignalSvc: svc}
// }
// func (ctrl *CustomerSignalController) GetTraderSignalsForCustomer(c *gin.Context) {
// 	customerID := c.MustGet("userID").(uint)
// 	traderIDStr := c.Param("trader_id")
// 	traderID, err := strconv.ParseUint(traderIDStr, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid trader ID"})
// 		return
// 	}

// 	signals, err := ctrl.CustomerSignalSvc.GetTraderSignalsForCustomer(c.Request.Context(), customerID, uint(traderID))
// 	if err != nil {
// 		statusCode := http.StatusInternalServerError
// 		if errors.Is(err, service.ErrNotSubscribed) {
// 			statusCode = http.StatusForbidden // Customer not subscribed
// 		}
// 		c.JSON(statusCode, gin.H{"message": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, signals)
// }
// func (ctrl *CustomerSignalController) GetSignalCardForCustomer(c *gin.Context) {
// 	customerID := c.MustGet("userID").(uint)
// 	traderIDStr := c.Param("trader_id")
// 	traderID, err := strconv.ParseUint(traderIDStr, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid trader ID"})
// 		return
// 	}

// 	signalIDStr := c.Param("signal_id")
// 	signalID, err := strconv.ParseUint(signalIDStr, 10, 32)
// 	if err != nil {
// 		c.JSON(http.StatusBadRequest, gin.H{"message": "Invalid signal ID"})
// 		return
// 	}

// 	signal, err := ctrl.CustomerSignalSvc.GetSignalCardForCustomer(c.Request.Context(), customerID, uint(traderID), uint(signalID))
// 	if err != nil {
// 		statusCode := http.StatusInternalServerError
// 		if errors.Is(err, service.ErrNotSubscribed) {
// 			statusCode = http.StatusForbidden
// 		} else if errors.Is(err, errors.New("signal does not belong to the specified trader")) {
// 			statusCode = http.StatusForbidden // Or 404 if you want to hide its existence
// 		} else if errors.Is(err, gorm.ErrRecordNotFound) { // Assuming repo returns this
// 			statusCode = http.StatusNotFound
// 		}
// 		c.JSON(statusCode, gin.H{"message": err.Error()})
// 		return
// 	}

// 	c.JSON(http.StatusOK, signal)
// }
