package routes

import (
	"hypha/api/internal/db"
	"hypha/api/internal/db/handlers"

	"github.com/gin-gonic/gin"
)

func InitResultsRoutes(router *gin.RouterGroup, dpOps db.DatabaseOperations) {
	router.GET("/integration/:id", func(context *gin.Context) {
		handlers.GetResultsByIntegrationID(dpOps, context)
	})
	router.GET("/product/:productId", func(context *gin.Context) {
		handlers.GetResultsByProductID(dpOps, context)
	})
	router.POST("/results", func(context *gin.Context) {
		handlers.ReportResults(dpOps, context)
	})
}
