package routing

import (
	"github.com/CBE-Super-App/cbe-super-app-member-auth/internal/handlers/rest/external_call"
	"github.com/gin-gonic/gin"
)

// SetupExternalCallRoutes sets up the external call routes
func SetupExternalCallRoutes(router *gin.RouterGroup, handler *external_call.ExternalCallHandler) {
	externalCallGroup := router.Group("/external-calls")
	{
		// SMS endpoints
		externalCallGroup.POST("/sms/send", handler.SendSMS)

		// Generic external call endpoints
		externalCallGroup.POST("/call", handler.MakeExternalCall)

		// History and monitoring endpoints
		externalCallGroup.GET("/history", handler.GetExternalCallHistory)
		externalCallGroup.GET("/history/date-range", handler.GetExternalCallHistoryByDateRange)
		externalCallGroup.GET("/:id", handler.GetExternalCallById)
	}
}
