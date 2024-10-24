package routes

import (
	"hypha/api/internal/db"
	"hypha/api/internal/db/handlers"

	"github.com/gin-gonic/gin"
)

// InitProductRoutes initializes the product routes for the given router group.
// It sets up the POST and GET endpoints for creating and retrieving products and their integrations.
//
// Parameters:
// - router: The router group to which the routes will be added.
// - dbOps: The database operations interface used for database interactions.
//
// Routes:
// - POST /product: Calls CreateProduct to handle the creation of a new product.
// - GET /product/:id: Calls GetProduct to handle retrieving a product by ID.
// - GET /product/:id/integrations: Calls GetProductIntegrations to handle retrieving integrations for a product by ID.
// - GET /products: Calls GetAllProducts to handle retrieving all products.
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

// InitRelationshipRoutes initializes the relation routes for the given router group.
// It sets up the POST and GET endpoints for creating and retrieving relationships.
//
// Parameters:
// - router: The router group to which the routes will be added.
// - dbOps: The database operations interface used for database interactions.
//
// Routes:
// - POST /relationship: Calls CreateRelationship to handle the creation of a new relationship.
// - GET /relationship/:id: Calls GetRelationship to handle retrieving a relationship by ID.
func InitRelationshipRoutes(router *gin.RouterGroup, dbOps db.DatabaseOperations) {
	router.POST("/relationship", func(context *gin.Context) {
		handlers.CreateRelationship(dbOps, context)
	})
	router.GET("/relationship/:id", func(context *gin.Context) {
		handlers.GetRelationship(dbOps, context)
	})
}

// InitRuleRoutes initializes the rule routes for the given router group.
// It sets up the POST and GET endpoints for creating and retrieving results rules.
//
// Parameters:
// - router: The router group to which the routes will be added.
// - dbOps: The database operations interface used for database interactions.
//
// Routes:
// - POST /results-rule: Calls CreateResultsRule to handle the creation of a new results rule.
// - GET /results-rule/:id: Calls GetResultsRule to handle retrieving a results rule by ID.
// - GET /results-rule/relation/:id: Calls GetRulesByRelationIDto handle retrieving results rules using a relation ID
func InitRuleRoutes(router *gin.RouterGroup, dbOps db.DatabaseOperations) {
	router.POST("/results-rule", func(context *gin.Context) {
		handlers.CreateResultsRule(dbOps, context)
	})
	router.GET("/results-rule/:id", func(context *gin.Context) {
		handlers.GetResultsRule(dbOps, context)
	})
	router.GET("/results-rule/relation/:id", func(context *gin.Context) {
		handlers.GetRulesByRelationID(dbOps, context)
	})
}
