package controllers

import (
	"net/http"

	"github.com/fathimasithara01/tradeverse/internal/admin/service"
	"github.com/gin-gonic/gin"
)

type WebConfigurationController struct {
	webConfigService service.WebConfigurationService
}

func NewWebConfigurationController(webConfigService service.WebConfigurationService) *WebConfigurationController {
	return &WebConfigurationController{webConfigService: webConfigService}
}

func (ctrl *WebConfigurationController) GetWebConfigurationPage(c *gin.Context) {
	config, err := ctrl.webConfigService.GetWebConfiguration()
	if err != nil {
		c.HTML(http.StatusInternalServerError, "admin_web_configuration.html", gin.H{
			"error": "Failed to load web configuration",
		})
		return
	}

	countries := []string{"United Arab Emirates", "United States", "India", "United Kingdom", "Canada"}
	currencies := []string{"United Arab Emirates Dirham (AED)", "United States Dollar (USD)", "Indian Rupee (INR)", "British Pound (GBP)", "Canadian Dollar (CAD)"}
	timezones := []string{"Asia/Dubai", "America/New_York", "Asia/Kolkata", "Europe/London", "America/Toronto"}

	c.HTML(http.StatusOK, "admin_web_configuration.html", gin.H{
		"config":     config,
		"success":    c.Query("success"), 
		"error":      c.Query("error"),   
		"countries":  countries,
		"currencies": currencies,
		"timezones":  timezones,
	})
}

func (ctrl *WebConfigurationController) UpdateWebConfiguration(c *gin.Context) {
	primaryCountry := c.PostForm("primary_country")
	primaryCurrency := c.PostForm("primary_currency")
	primaryTimezone := c.PostForm("primary_timezone")
	// Removed: filesystemConfig := c.PostForm("filesystem_config") 

	if primaryCountry == "" || primaryCurrency == "" || primaryTimezone == "" {
		c.Redirect(http.StatusFound, "/admin/web-configuration?error=All fields are required")
		return
	}

	err := ctrl.webConfigService.UpdateWebConfiguration(primaryCountry, primaryCurrency, primaryTimezone)
	if err != nil {
		c.Redirect(http.StatusFound, "/admin/web-configuration?error=Failed to update configuration")
		return
	}

	c.Redirect(http.StatusFound, "/admin/web-configuration?success=Configuration updated successfully")
}
