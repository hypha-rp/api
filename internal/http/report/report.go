package report

import (
	"bytes"
	"encoding/xml"
	"hypha/api/internal/db"
	"io"
	"net/http"

	"hypha/api/internal/utils/logging"
	"hypha/api/internal/utils/results"

	"github.com/gin-gonic/gin"
)

var logger = logging.Logger

// InitReportRoutes initializes the report routes for the given router group.
// It sets up the POST route for reporting results.
//
// Parameters:
// - router: The router group to which the routes will be added.
// - dbOperations: The database operations interface for interacting with the database.
func InitReportRoutes(router *gin.RouterGroup, dbOperations db.DatabaseOperations) {
	router.POST("/results", func(c *gin.Context) {
		ReportResults(c, dbOperations)
	})
}

// ReportResults handles the reporting of test results.
// It processes the uploaded JUnit XML file, parses the results, and stores them in the database.
//
// Parameters:
// - c: The Gin context for the current request.
// - dbOperations: The database operations interface for interacting with the database.
func ReportResults(c *gin.Context, dbOperations db.DatabaseOperations) {
	var junitTestSuites results.JUnitTestSuites
	var product db.Product

	productId := c.PostForm("productId")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "productId is required"})
		return
	}

	if dbOperations == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database operations not initialized"})
		return
	}

	if err := dbOperations.First(&product, "id = ?", productId); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid productId"})
		return
	}

	file, err := c.FormFile("file")
	if err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "File upload failed"})
		return
	}

	fileContent, err := file.Open()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to open file"})
		return
	}
	defer fileContent.Close()

	byteValue, err := io.ReadAll(fileContent)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read file"})
		return
	}

	if !containsTestsuitesTag(byteValue) {
		byteValue = wrapInTestsuitesTag(byteValue)
	}

	if err := xml.Unmarshal(byteValue, &junitTestSuites); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid XML format. Not JUnit"})
		logger.Error().Err(err).Msg("Failed to unmarshal JUnit XML")
		return
	}

	if err := results.ParseJUnitResults(junitTestSuites, dbOperations, productId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

// containsTestsuitesTag checks if the given XML content contains a <testsuites> tag.
// Parameters:
// - xmlContent: The XML content to check as a byte slice.
// Returns:
// - A boolean indicating whether the <testsuites> tag is present in the XML content.
func containsTestsuitesTag(xmlContent []byte) bool {
	return bytes.Contains(xmlContent, []byte("<testsuites"))
}

// wrapInTestsuitesTag wraps the given XML content in a <testsuites> tag.
// Parameters:
// - xmlContent: The XML content to wrap as a byte slice.
// Returns:
// - A new byte slice with the XML content wrapped in a <testsuites> tag.
func wrapInTestsuitesTag(xmlContent []byte) []byte {
	return append([]byte("<testsuites>"), append(xmlContent, []byte("</testsuites>")...)...)
}
