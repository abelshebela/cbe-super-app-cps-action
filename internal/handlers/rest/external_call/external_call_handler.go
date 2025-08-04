package external_call

import (
	"net/http"
	"time"

	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/constants/dto"
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/service/external_call"
	"github.com/gin-gonic/gin"
	"gitlab.com/bersufekadgetachew/cbe-super-app-shared/shared/utils"
)

type ExternalCallHandler struct {
	externalCallService *external_call.ExternalCallService
	logger              utils.Logger
}

func NewExternalCallHandler(externalCallService *external_call.ExternalCallService, logger utils.Logger) *ExternalCallHandler {
	return &ExternalCallHandler{
		externalCallService: externalCallService,
		logger:              logger,
	}
}

// SendSMS handles SMS sending requests
func (h *ExternalCallHandler) SendSMS(c *gin.Context) {
	var request dto.SMSRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Errorf("Failed to bind SMS request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	response, err := h.externalCallService.SendSMS(c.Request.Context(), &request)
	if err != nil {
		h.logger.Errorf("Failed to send SMS: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to send SMS",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": response.Success,
		"message": response.Message,
		"data":    response.Data,
	})
}

// MakeExternalCall handles generic external API calls
func (h *ExternalCallHandler) MakeExternalCall(c *gin.Context) {
	var request dto.ExternalCallRequest
	if err := c.ShouldBindJSON(&request); err != nil {
		h.logger.Errorf("Failed to bind external call request: %v", err)
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid request format",
			"error":   err.Error(),
		})
		return
	}

	response, err := h.externalCallService.MakeGenericExternalCall(c.Request.Context(), &request)
	if err != nil {
		h.logger.Errorf("Failed to make external call: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to make external call",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": response.StatusCode >= 200 && response.StatusCode < 300,
		"data":    response,
	})
}

// GetExternalCallHistory retrieves external call history by URL
func (h *ExternalCallHandler) GetExternalCallHistory(c *gin.Context) {
	url := c.Query("url")
	if url == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "URL parameter is required",
		})
		return
	}

	history, err := h.externalCallService.GetExternalCallHistory(c.Request.Context(), url)
	if err != nil {
		h.logger.Errorf("Failed to get external call history: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve external call history",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    history,
		"count":   len(history),
	})
}

// GetExternalCallHistoryByDateRange retrieves external call history within a date range
func (h *ExternalCallHandler) GetExternalCallHistoryByDateRange(c *gin.Context) {
	startDateStr := c.Query("start_date")
	endDateStr := c.Query("end_date")

	if startDateStr == "" || endDateStr == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Both start_date and end_date parameters are required",
		})
		return
	}

	startDate, err := time.Parse("2006-01-02", startDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid start_date format. Use YYYY-MM-DD",
			"error":   err.Error(),
		})
		return
	}

	endDate, err := time.Parse("2006-01-02", endDateStr)
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "Invalid end_date format. Use YYYY-MM-DD",
			"error":   err.Error(),
		})
		return
	}

	// Set end date to end of day
	endDate = endDate.Add(23*time.Hour + 59*time.Minute + 59*time.Second)

	history, err := h.externalCallService.GetExternalCallHistoryByDateRange(c.Request.Context(), startDate, endDate)
	if err != nil {
		h.logger.Errorf("Failed to get external call history by date range: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve external call history",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    history,
		"count":   len(history),
	})
}

// GetExternalCallById retrieves a specific external call by ID
func (h *ExternalCallHandler) GetExternalCallById(c *gin.Context) {
	id := c.Param("id")
	if id == "" {
		c.JSON(http.StatusBadRequest, gin.H{
			"success": false,
			"message": "ID parameter is required",
		})
		return
	}

	externalCall, err := h.externalCallService.GetExternalCallById(c.Request.Context(), id)
	if err != nil {
		h.logger.Errorf("Failed to get external call by ID: %v", err)
		c.JSON(http.StatusInternalServerError, gin.H{
			"success": false,
			"message": "Failed to retrieve external call",
			"error":   err.Error(),
		})
		return
	}

	c.JSON(http.StatusOK, gin.H{
		"success": true,
		"data":    externalCall,
	})
}
