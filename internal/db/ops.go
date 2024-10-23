package db

import (
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-orm/gorm"
	"github.com/google/uuid"
)

// DatabaseOperations defines the interface for database operations.
// It includes methods for obtaining a database connection, creating a record, and fetching the first record that matches the criteria.
type DatabaseOperations interface {
	Connection() *gorm.DB
	Create(value interface{}) error
	First(out interface{}, where ...interface{}) error
}

// CreateResource handles the creation of a new resource in the database.
// It binds the JSON request body to the resource, validates it, and creates the resource in the database.
//
// Parameters:
// - dbOps: The database operations interface used for database interactions.
// - context: The Gin context that provides request and response handling.
// - resource: The resource object to be created.
//
// Responses:
// - 400 Bad Request: If the request body is invalid or cannot be bound to the resource.
// - 500 Internal Server Error: If there is an error creating the resource in the database.
// - 200 OK: If the resource is successfully created, returns the created resource.
func CreateResource(dbOps DatabaseOperations, context *gin.Context, resource interface{}) {
	if err := context.ShouldBindJSON(resource); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": err.Error()})
		return
	}

	if err := dbOps.Create(resource); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	context.JSON(http.StatusOK, resource)
}

// GetResource retrieves an existing resource by its ID.
// It fetches the resource from the database and returns it in the response.
//
// Parameters:
// - dbOps: The database operations interface used for database interactions.
// - context: The Gin context that provides request and response handling.
// - resource: The resource object to be retrieved.
// - idParam: The name of the URL parameter that contains the resource ID.
// - resourceName: The name of the resource, used for error messages.
//
// Responses:
// - 404 Not Found: If the resource is not found in the database.
// - 200 OK: If the resource is successfully retrieved, returns the resource object.
func GetResource(dbOps DatabaseOperations, context *gin.Context, resource interface{}, idParam string, resourceName string) {
	resourceID := context.Param(idParam)

	if err := dbOps.First(resource, "id = ?", resourceID); err != nil {
		context.JSON(http.StatusNotFound, gin.H{"error": resourceName + " not found"})
		return
	}

	context.JSON(http.StatusOK, resource)
}

// GenerateUniqueID generates a unique identifier using UUID.
//
// Returns:
// - string: A unique identifier string.
func GenerateUniqueID() string {
	return uuid.New().String()
}
