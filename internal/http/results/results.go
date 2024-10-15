package results

import (
	"hypha/api/internal/db"
	"hypha/api/internal/utils/logging"
	"net/http"

	"github.com/gin-gonic/gin"
)

var log = logging.Logger

// InitResultsRoutes initializes the results routes for the given router group.
// It sets up the GET routes for retrieving results by integration ID and product ID.
//
// Parameters:
// - router: The router group to which the routes will be added.
// - dbOperations: The database operations interface for interacting with the database.
func InitResultsRoutes(router *gin.RouterGroup, dbOperations db.DatabaseOperations) {
	router.GET("/results/integration/:id", func(context *gin.Context) {
		GetResultsByIntegrationID(dbOperations, context)
	})
	router.GET("/results/product/:productId", func(c *gin.Context) {
		GetResultsByProductID(c, dbOperations)
	})
}

// GetResultsByIntegrationID retrieves test results based on the integration ID.
// It fetches the test suite and test case IDs, retrieves the test suites, filters the test cases,
// and fetches the results and associated products from the database.
//
// Parameters:
// - dbOperations: The database operations interface for interacting with the database.
// - context: The Gin context for the current request.
func GetResultsByIntegrationID(dbOperations db.DatabaseOperations, context *gin.Context) {
	integrationID := context.Param("id")
	if integrationID == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "integration ID is required"})
		return
	}

	db := dbOperations.Connection()

	testSuiteIDs, testCaseIDs, err := getTestSuiteAndCaseIDs(db, integrationID)
	if err != nil {
		logErrorAndRespond(context, "Database query error in getTestSuiteAndCaseIDs", err)
		return
	}

	if len(testSuiteIDs) == 0 && len(testCaseIDs) == 0 {
		context.JSON(http.StatusNotFound, gin.H{"error": "No test suites or test cases found for the given integration ID"})
		return
	}

	testSuites, err := getTestSuites(db, testSuiteIDs, testCaseIDs)
	if err != nil {
		logErrorAndRespond(context, "Database query error in getTestSuites", err)
		return
	}

	filterTestCases(testSuites, integrationID)

	resultMap := createResultMap(testSuites)

	results, err := fetchResultsAndProducts(db, resultMap)
	if err != nil {
		logErrorAndRespond(context, "Database query error in fetchResultsAndProducts", err)
		return
	}

	context.JSON(http.StatusOK, results)
}

// GetResultsByProductID retrieves test results based on the product ID.
// It fetches the results and associated test suites, test cases, and properties from the database.
//
// Parameters:
// - c: The Gin context for the current request.
// - dbOperations: The database operations interface for interacting with the database.
func GetResultsByProductID(c *gin.Context, dbOperations db.DatabaseOperations) {
	productId := c.Param("productId")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "productId is required"})
		return
	}

	var results []db.Result

	db := dbOperations.Connection()

	if err := db.Where("product_id = ?", productId).
		Preload("TestSuites").
		Preload("TestSuites.TestCases").
		Preload("TestSuites.Properties").
		Preload("TestSuites.TestCases.Properties").
		Find(&results).Error; err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to retrieve results"})
		return
	}

	c.JSON(http.StatusOK, results)
}
