package handlers

import (
	"bytes"
	"encoding/json"
	"hypha/api/internal/db"
	"hypha/api/internal/db/tables"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

// CreateRelationship handles the creation of a new relationship between two products.
// It reads the request body, validates the input, checks for existing relationships, and creates a new relationship record in the database.
//
// Parameters:
// - dbOps: The database operations interface used for database interactions.
// - context: The Gin context that provides request and response handling.
//
// Request Body:
// The request body should be a JSON object containing the following fields:
//   - objectID1 (string): The ID of the first object.
//   - objectID2 (string): The ID of the second object.
//   - relationshipType (string): The type of the relationship (e.g., "integration", "dependency").
//
// Responses:
// - 400 Bad Request: If the request body is invalid, if the object IDs are empty or the same, or if a relationship already exists for the given objects.
// - 201 Created: If the relationship is successfully created.
func CreateRelationship(dbOps db.DatabaseOperations, context *gin.Context) {
	// Read the request body and parse the JSON
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

	objectID1 := requestBody["objectID1"]
	objectID2 := requestBody["objectID2"]
	relationshipType := requestBody["relationshipType"]

	if objectID1 == "" || objectID2 == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Object IDs cannot be empty"})
		return
	}

	if objectID1 == objectID2 {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Object IDs cannot be the same"})
		return
	}

	// Check if the relationship already exists
	var existingRelationship tables.Relationship
	if err := dbOps.Connection().
		Where("object_ids @> ARRAY[?]::text[] AND object_ids @> ARRAY[?]::text[] AND relationship_type = ?", objectID1, objectID2, relationshipType).
		First(&existingRelationship).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Relationship already exists"})
		return
	}

	// Create the new relationship
	relationship := tables.Relationship{
		ID:               db.GenerateUniqueID(),
		ObjectIDs:        pq.StringArray{objectID1, objectID2},
		RelationshipType: relationshipType,
	}

	if err := dbOps.Connection().Create(&relationship).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	context.JSON(http.StatusCreated, relationship)
}

// GetRelationship retrieves an existing relationship by its ID.
// It fetches the relationship from the database, including related products, and returns it in the response.
//
// Parameters:
// - dbOps: The database operations interface used for database interactions.
// - context: The Gin context that provides request and response handling.
//
// Path Parameters:
// - id (string): The ID of the relationship to retrieve.
//
// Responses:
// - 500 Internal Server Error: If there is an error retrieving the relationship or related products from the database.
// - 200 OK: If the relationship is successfully retrieved, returns the relationship object along with the related products.
func GetRelationship(dbOps db.DatabaseOperations, context *gin.Context) {
	var existingRelationship tables.Relationship
	if err := dbOps.Connection().
		Where("id = ?", context.Param("id")).
		First(&existingRelationship).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	var objects []tables.Product
	for _, objectID := range existingRelationship.ObjectIDs {
		var product tables.Product
		if err := dbOps.Connection().
			Where("id::text = ?", objectID).
			First(&product).Error; err != nil {
			context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
			return
		}
		objects = append(objects, product)
	}

	response := gin.H{
		"id":               existingRelationship.ID,
		"objectIDs":        existingRelationship.ObjectIDs,
		"relationshipType": existingRelationship.RelationshipType,
		"objects":          objects,
	}

	context.JSON(http.StatusOK, response)
}
