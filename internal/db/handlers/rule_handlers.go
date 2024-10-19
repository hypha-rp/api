package handlers

import (
	"hypha/api/internal/db"

	"github.com/gin-gonic/gin"
)

// CreateResultsRule handles the creation of a new results rule.
// It generates a unique ID for the new rule and creates the rule in the database.
//
// Parameters:
// - dbOps: The database operations interface used for database interactions.
// - context: The Gin context that provides request and response handling.
//
// Request Body:
// The request body should be a JSON object containing the fields required for a ResultsRule.
//
// Responses:
// - 201 Created: If the results rule is successfully created.
func CreateResultsRule(dbOps db.DatabaseOperations, context *gin.Context) {
	var newRule db.ResultsRule
	newRule.ID = db.GenerateUniqueID()
	db.CreateResource(dbOps, context, &newRule)
}

// GetResultsRule retrieves an existing results rule by its ID.
// It fetches the results rule from the database and returns it in the response.
//
// Parameters:
// - dbOps: The database operations interface used for database interactions.
// - context: The Gin context that provides request and response handling.
//
// Path Parameters:
// - id (string): The ID of the results rule to retrieve.
//
// Responses:
// - 200 OK: If the results rule is successfully retrieved, returns the results rule object.
func GetResultsRule(dbOps db.DatabaseOperations, context *gin.Context) {
	var existingRule db.ResultsRule
	db.GetResource(dbOps, context, &existingRule, "id", "ResultsRule")
}

// GetResultsRuleByRelationID retrieves a results rule by its relation ID.
// It fetches the results rule from the database based on the provided relation ID.
//
// Parameters:
// - dbOps: The database operations interface used for database interactions.
// - relationID: The relation ID to filter by.
//
// Returns:
// - *db.ResultsRule: The results rule object if found.
// - error: An error if any database operation fails.
func GetResultsRuleByRelationID(dbOps db.DatabaseOperations, relationID string) (*db.ResultsRule, error) {
	var rule db.ResultsRule
	err := dbOps.Connection().Where("relation_id = ?", relationID).First(&rule).Error
	if err != nil {
		return nil, err
	}
	return &rule, nil
}
