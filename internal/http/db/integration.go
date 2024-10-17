package db

import (
	"bytes"
	"encoding/json"
	"hypha/api/internal/db"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// InitIntegrationRoutes initializes the integration routes for the given router group.
// It sets up the POST and GET endpoints for creating and retrieving integrations.
//
// Parameters:
//
//	router (*gin.RouterGroup): The router group to which the routes will be added.
//	dbOperations (db.DatabaseOperations): The database operations interface used for database interactions.
//
// Routes:
//
//	POST /integration: Calls CreateIntegration to handle the creation of a new integration.
//	GET /integration/:id: Calls GetIntegration to handle retrieving an integration by ID.
func InitIntegrationRoutes(router *gin.RouterGroup, dbOperations db.DatabaseOperations) {
	router.POST("/integration", func(context *gin.Context) {
		CreateIntegration(dbOperations, context)
	})
	router.GET("/integration/:id", func(context *gin.Context) {
		GetIntegration(dbOperations, context)
	})
}

// CreateIntegration handles the creation of a new integration between two products.
// It reads the request body, validates the input, checks for existing integrations, and creates a new integration record in the database.
//
// Parameters:
//
//	dbOperations (db.DatabaseOperations): The database operations interface used for database interactions.
//	context (*gin.Context): The Gin context that provides request and response handling.
//
// Request Body:
//
//	The request body should be a JSON object containing the following fields:
//	  - productID1 (string): The ID of the first product.
//	  - productID2 (string): The ID of the second product.
//
// Responses:
//
//	400 Bad Request: If the request body is invalid, if the product IDs are empty or the same, or if an integration already exists for the given products.
//	201 Created: If the integration is successfully created.
func CreateIntegration(dbOperations db.DatabaseOperations, context *gin.Context) {
	var requestBody map[string]string
	body, err := io.ReadAll(context.Request.Body)
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Failed to read request body"})
		return
	}

	context.Request.Body = io.NopCloser(bytes.NewBuffer(body))

	if err := json.Unmarshal(body, &requestBody); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid JSON"})
		return
	}

	productID1 := requestBody["productID1"]
	productID2 := requestBody["productID2"]

	if productID1 == "" || productID2 == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Product IDs cannot be empty"})
		return
	}

	if productID1 == productID2 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Cannot create integration for the same product"})
		return
	}

	var existingIntegration db.Integration
	if err := dbOperations.Connection().
		Where("(product_id1 = ? AND product_id2 = ?) OR (product_id1 = ? AND product_id2 = ?)", productID1, productID2, productID2, productID1).
		First(&existingIntegration).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Integration already exists for the given products"})
		return
	}

	var newIntegration db.Integration
	newIntegration.ID = db.GenerateUniqueID()
	db.CreateResource(dbOperations, context, &newIntegration)
}

// GetIntegration retrieves an existing integration by its ID.
// It fetches the integration from the database, including related products, and returns it in the response.
//
// Parameters:
//
//	dbOperations (db.DatabaseOperations): The database operations interface used for database interactions.
//	context (*gin.Context): The Gin context that provides request and response handling.
//
// Path Parameters:
//
//	id (string): The ID of the integration to retrieve.
//
// Responses:
//
//	500 Internal Server Error: If there is an error retrieving the integration from the database.
//	200 OK: If the integration is successfully retrieved, returns the integration object.
func GetIntegration(dbOperations db.DatabaseOperations, context *gin.Context) {
	var existingIntegration db.Integration
	if err := dbOperations.Connection().
		Preload("Product1").
		Preload("Product2").
		Where("id = ?", context.Param("id")).
		First(&existingIntegration).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	context.JSON(http.StatusOK, existingIntegration)
}
