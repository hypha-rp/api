package handlers

import (
	"encoding/xml"
	"hypha/api/internal/db"
	"hypha/api/internal/db/tables"
	"hypha/api/internal/utils/db/queries"
	"hypha/api/internal/utils/results"
	"io"
	"net/http"

	"github.com/gin-gonic/gin"
)

// GetResultsByRelationID retrieves test results based on the relation ID.
// It fetches the test suite and test case IDs, retrieves the test suites, filters the test cases,
// and fetches the results and associated products from the database.
//
// Parameters:
// - dbOps: The database operations interface for interacting with the database.
// - context: The Gin context for the current request.
func GetResultsByRelationID(dbOps db.DatabaseOperations, context *gin.Context) {
	relationID := context.Param("id")
	db := dbOps.Connection()

	if relationID == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "relation ID is required"})
		return
	}

	rules, err := queries.FetchRulesByRelationID(dbOps, relationID)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	results, err := queries.FetchResultsByRules(db, rules)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	context.JSON(http.StatusOK, results)
}

// GetResultsByProductID retrieves test results based on the product ID.
// It fetches the results and associated test suites, test cases, and properties from the database.
//
// Parameters:
// - dbOps: The database operations interface for interacting with the database.
// - context: The Gin context for the current request.
func GetResultsByProductID(dbOps db.DatabaseOperations, context *gin.Context) {
	productId := context.Param("productId")
	if productId == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "productId is required"})
		return
	}

	var results []tables.Result

	db := dbOps.Connection()

	if err := db.Where("product_id = ?", productId).
		Preload("TestSuites").
		Preload("TestSuites.TestCases").
		Preload("TestSuites.Properties").
		Preload("TestSuites.TestCases.Properties").
		Find(&results).Error; err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	context.JSON(http.StatusOK, results)
}

// ReportResults handles the reporting of test results.
// It processes the uploaded JUnit XML file, parses the results, and stores them in the database.
//
// Parameters:
// - dbOps: The database operations interface for interacting with the database.
// - context: The Gin context for the current request.
func ReportResults(dpOps db.DatabaseOperations, context *gin.Context) {
	var junitTestSuites results.JUnitTestSuites
	var product tables.Product

	productId := context.PostForm("productId")
	if productId == "" {
		context.JSON(http.StatusBadRequest, gin.H{"error": "productId is required"})
		return
	}

	if dpOps == nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if err := dpOps.First(&product, "id = ?", productId); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid productId"})
		return
	}

	file, err := context.FormFile("file")
	if err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	fileContent, err := file.Open()
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}
	defer fileContent.Close()

	byteValue, err := io.ReadAll(fileContent)
	if err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	if !results.ContainsTestsuitesTag(byteValue) {
		byteValue = results.WrapInTestsuitesTag(byteValue)
	}

	if err := xml.Unmarshal(byteValue, &junitTestSuites); err != nil {
		context.JSON(http.StatusBadRequest, gin.H{"error": "Invalid XML format. Not JUnit"})
		return
	}

	if err := results.ParseJUnitResults(junitTestSuites, dpOps, productId); err != nil {
		context.JSON(http.StatusInternalServerError, gin.H{"error": "Internal server error"})
		return
	}

	context.JSON(http.StatusOK, gin.H{"status": "success"})
}
