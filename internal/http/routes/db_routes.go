package routes

import (
	"hypha/api/internal/db"
	"hypha/api/internal/db/handlers"

	"github.com/gin-gonic/gin"
)

func InitProductRoutes(router *gin.RouterGroup, dbOps db.DatabaseOperations) {
	router.POST("/product", func(context *gin.Context) {
		handlers.CreateProduct(dbOps, context)
	})
	router.GET("/product/:id", func(context *gin.Context) {
		handlers.GetProduct(dbOps, context)
	})
	router.GET("/product/:id/integrations", func(context *gin.Context) {
		handlers.GetProductIntegrations(dbOps, context)
	})
	router.GET("/products", func(context *gin.Context) {
		handlers.GetAllProducts(dbOps, context)
	})
}

// InitIntegrationRoutes initializes the integration routes for the given router group.
// It sets up the POST and GET endpoints for creating and retrieving integrations.
//
// Parameters:
//
//	router (*gin.RouterGroup): The router group to which the routes will be added.
//	dbOps (db.DatabaseOperations): The database operations interface used for database interactions.
//
// Routes:
//
//	POST /integration: Calls CreateIntegration to handle the creation of a new integration.
//	GET /integration/:id: Calls GetIntegration to handle retrieving an integration by ID.
func InitRelationRoutes(router *gin.RouterGroup, dbOps db.DatabaseOperations) {
	router.POST("/relation", func(context *gin.Context) {
		handlers.CreateRelationship(dbOps, context)
	})
	router.GET("/relation/:id", func(context *gin.Context) {
		handlers.GetRelationship(dbOps, context)
	})
}

func InitRuleRoutes(router *gin.RouterGroup, dbOps db.DatabaseOperations) {
	router.POST("/results-rule", func(context *gin.Context) {
		handlers.CreateResultsRule(dbOps, context)
	})
	router.GET("/results-rule/:id", func(context *gin.Context) {
		handlers.GetResultsRule(dbOps, context)
	})
}
