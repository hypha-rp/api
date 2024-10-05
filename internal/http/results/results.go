package results

import (
	"hypha/api/internal/db/ops"
	"hypha/api/internal/utils/logging"
	"net/http"

	"github.com/gin-gonic/gin"
)

var log = logging.Logger

func InitResultsRoutes(router *gin.RouterGroup, dbOperations ops.DatabaseOperations) {
	router.GET("/results/integration/:id", func(context *gin.Context) {
		GetResultsByIntegrationID(dbOperations, context)
	})
}

func GetResultsByIntegrationID(dbOperations ops.DatabaseOperations, context *gin.Context) {
	integrationID := context.Param("id")
	if integrationID == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "integration ID is required"})
		return
	}

	db := dbOperations.Connection()

	testSuiteIDs, testCaseIDs, err := getTestSuiteAndCaseIDs(db, integrationID)
	if err != nil {
		log.Error().Msgf("Database query error in getTestSuiteAndCaseIDs: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "There was a problem processing your request"})
		return
	}

	if len(testSuiteIDs) == 0 && len(testCaseIDs) == 0 {
		context.JSON(http.StatusNotFound, gin.H{"error": "No test suites or test cases found for the given integration ID"})
		return
	}

	testSuites, err := getTestSuites(db, testSuiteIDs, testCaseIDs)
	if err != nil {
		log.Error().Msgf("Database query error in getTestSuites: %v", err)
		context.JSON(http.StatusInternalServerError, gin.H{"error": "There was a problem processing your request"})
		return
	}

	filterTestCases(testSuites, integrationID)

	context.JSON(http.StatusOK, testSuites)
}
