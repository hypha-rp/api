package routes

import (
	"hypha/api/internal/db"
	"hypha/api/internal/db/handlers"

	"github.com/gin-gonic/gin"
)

// InitResultsRoutes initializes the results routes for the given router group.
// It sets up the GET and POST endpoints for retrieving and reporting results.
//
// Parameters:
// - router: The router group to which the routes will be added.
// - dpOps: The database operations interface used for database interactions.
//
// Routes:
// - GET /integration/:id: Calls GetResultsByIntegrationID to handle retrieving results by integration ID.
// - GET /product/:productId: Calls GetResultsByProductID to handle retrieving results by product ID.
// - POST /results: Calls ReportResults to handle reporting new results.
func InitResultsRoutes(router *gin.RouterGroup, dpOps db.DatabaseOperations) {
	router.GET("/relationship/:id", func(context *gin.Context) {
		handlers.GetResultsByRelationID(dpOps, context)
	})
	router.GET("/product/:productId", func(context *gin.Context) {
		handlers.GetResultsByProductID(dpOps, context)
	})
	router.POST("/", func(context *gin.Context) {
		handlers.ReportResults(dpOps, context)
	})
}
