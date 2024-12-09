package routes

import (
	"book_service/pkg/routes/api"
	"book_service/pkg/routes/common"
	"github.com/gin-gonic/gin"
)

func RegisterRoutes(router *gin.Engine) {
	api.ApiRouter(router)
	common.HealthRoutes(router)
}
