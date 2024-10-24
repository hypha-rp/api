package handlers

import (
	"hypha/api/internal/db"
	"hypha/api/internal/db/tables"
	"net/http"

	"github.com/gin-gonic/gin"
)

// CreateProduct handles the creation of a new product.
// It generates a unique ID for the new product and creates the product in the database.
//
// Parameters:
// - dbOps: The database operations interface used for database interactions.
// - context: The Gin context that provides request and response handling.
//
// Request Body:
// The request body should be a JSON object containing the fields required for a Product.
//
// Responses:
// - 201 Created: If the product is successfully created.
func CreateProduct(dbOps db.DatabaseOperations, context *gin.Context) {
	var newProduct tables.Product
	newProduct.ID = db.GenerateUniqueID()
	db.CreateResource(dbOps, context, &newProduct)
}

// GetProduct retrieves an existing product by its ID.
// It fetches the product from the database and returns it in the response.
//
// Parameters:
// - dbOps: The database operations interface used for database interactions.
// - context: The Gin context that provides request and response handling.
//
// Path Parameters:
// - id (string): The ID of the product to retrieve.
//
// Responses:
// - 200 OK: If the product is successfully retrieved, returns the product object.
func GetProduct(dbOps db.DatabaseOperations, context *gin.Context) {
	var existingProduct tables.Product
	db.GetResource(dbOps, context, &existingProduct, "id", "Product")
}

// GetProductIntegrations retrieves integrations for a product by its ID.
// It fetches the integrations from the database and returns them in the response.
//
// Parameters:
// - dbOps: The database operations interface used for database interactions.
// - context: The Gin context that provides request and response handling.
//
// Path Parameters:
// - id (string): The ID of the product to retrieve integrations for.
//
// Responses:
// - 200 OK: If the integrations are successfully retrieved, returns the integrations object.
func GetProductIntegrations(dbOps db.DatabaseOperations, context *gin.Context) {
	var integrations []tables.Relationship
	productID := context.Param("id")
	if err := dbOps.Connection().
		Where("relationship_type = ? AND ? = ANY(object_ids)", "integration", productID).
		Find(&integrations).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	for i, integration := range integrations {
		var objects []tables.ObjectInterface
		for _, objectID := range integration.ObjectIDs {
			var product tables.Product
			if err := dbOps.Connection().Where("id = ?", objectID).First(&product).Error; err == nil {
				objects = append(objects, product)
				continue
			}
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		integrations[i].Objects = objects
	}
	context.JSON(http.StatusOK, integrations)
}

// GetAllProducts retrieves all products, optionally filtered by name.
// It fetches the products from the database and returns them in the response.
//
// Parameters:
// - dbOps: The database operations interface used for database interactions.
// - context: The Gin context that provides request and response handling.
//
// Query Parameters:
// - name (string): Optional. The name to filter products by.
//
// Responses:
// - 200 OK: If the products are successfully retrieved, returns the products object.
func GetAllProducts(dbOps db.DatabaseOperations, context *gin.Context) {
	var products []tables.Product
	name := context.Query("name")
	query := dbOps.Connection()
	if name != "" {
		query = query.Where("full_name ILIKE ? OR short_name ILIKE ?", "%"+name+"%", "%"+name+"%")
	}
	if err := query.Find(&products).Error; err != nil {
		context.JSON(500, gin.H{"error": err.Error()})
		return
	}
	context.JSON(200, products)
}
