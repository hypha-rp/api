package handlers

import (
	"bytes"
	"encoding/json"
	"fmt"
	"hypha/api/internal/db"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
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
		context.JSON(http.StatusBadRequest, gin.H{"error": "Cannot create relationship for the same object"})
		return
	}

	var existingRelationship db.Relationship
	if err := dbOps.Connection().
		Where("object_ids @> ARRAY[?]::uuid[] AND object_ids @> ARRAY[?]::uuid[] AND relationship_type = ?", objectID1, objectID2, relationshipType).
		First(&existingRelationship).Error; err == nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": fmt.Sprintf("A %s relationship already exists for the given objects", relationshipType)})
		return
	}

	var newRelationship db.Relationship
	newRelationship.RelationID = db.GenerateUniqueID()
	newRelationship.ObjectIDs = []string{objectID1, objectID2}
	newRelationship.RelationshipType = relationshipType

	if err := dbOps.Connection().Create(&newRelationship).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	context.JSON(http.StatusCreated, newRelationship)
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
// - 500 Internal Server Error: If there is an error retrieving the relationship from the database.
// - 200 OK: If the relationship is successfully retrieved, returns the relationship object.
func GetRelationship(dbOps db.DatabaseOperations, context *gin.Context) {
	var existingRelationship db.Relationship
	if err := dbOps.Connection().
		Where("relation_id = ?", context.Param("id")).
		First(&existingRelationship).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	context.JSON(http.StatusOK, existingRelationship)
}
