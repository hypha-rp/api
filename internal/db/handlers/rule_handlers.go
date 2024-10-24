package handlers

import (
	"hypha/api/internal/db"
	"hypha/api/internal/db/tables"
	"hypha/api/internal/utils/db/queries"
	"hypha/api/internal/utils/logging"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/lib/pq"
)

var log = logging.Logger

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
	var requestBody struct {
		Expression string   `json:"expression"`
		AppliesTo  []string `json:"appliesTo"`
		RelationId string   `json:"relationId"`
	}

	if err := context.ShouldBindJSON(&requestBody); err != nil {
		log.Error().Err(err).Msg("Failed to bind JSON request body")
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid request body"})
		return
	}

	newRule := tables.ResultsRule{
		ID:             db.GenerateUniqueID(),
		Expression:     requestBody.Expression,
		AppliesTo:      pq.StringArray(requestBody.AppliesTo),
		RelationshipID: requestBody.RelationId,
		CreatedAt:      time.Now().UTC(),
		UpdatedAt:      time.Now().UTC(),
	}

	if err := dbOps.Create(&newRule); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to create results rule"})
		return
	}

	context.JSON(http.StatusCreated, gin.H{"message": "Results rule created successfully"})
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
	var existingRule tables.ResultsRule
	db.GetResource(dbOps, context, &existingRule, "id", "ResultsRule")
}

// GetRulesByRelationID retrieves results rules by their relation ID and returns them via the API.
//
// Parameters:
// - dbOps: The database operations interface used for database interactions.
// - context: The Gin context that provides request and response handling.
//
// Path Parameters:
// - id (string): The relation ID to filter by.
//
// Responses:
// - 200 OK: If the results rules are successfully retrieved, returns the results rule objects.
func GetRulesByRelationID(dbOps db.DatabaseOperations, context *gin.Context) {
	relationID := context.Param("id")
	rules, err := queries.FetchRulesByRelationID(dbOps, relationID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	context.JSON(http.StatusOK, rules)
}
