package report

import (
	"bytes"
	"encoding/xml"
	"hypha/api/internal/db/ops"
	"hypha/api/internal/db/tables"
	"io/ioutil"
	"net/http"

	"hypha/api/internal/utils/logging"
	"hypha/api/internal/utils/results/parse"
	"hypha/api/internal/utils/results/structs"

	"github.com/gin-gonic/gin"
)

var logger = logging.Logger

func InitReportRoutes(router *gin.RouterGroup, dbOperations ops.DatabaseOperations) {
	router.POST("/results", func(c *gin.Context) {
		ReportResults(c, dbOperations)
	})
	router.GET("/results/:productId", func(c *gin.Context) {
		GetResultsByProductID(c, dbOperations)
	})
}

func ReportResults(c *gin.Context, dbOperations ops.DatabaseOperations) {
	var junitTestSuites structs.JUnitTestSuites
	var product tables.Product

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

	byteValue, err := ioutil.ReadAll(fileContent)
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

	if err := parse.ParseJUnitResults(junitTestSuites, dbOperations, productId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}

func containsTestsuitesTag(xmlContent []byte) bool {
	return bytes.Contains(xmlContent, []byte("<testsuites>"))
}

func wrapInTestsuitesTag(xmlContent []byte) []byte {
	return append([]byte("<testsuites>"), append(xmlContent, []byte("</testsuites>")...)...)
}

func GetResultsByProductID(c *gin.Context, dbOperations ops.DatabaseOperations) {
	productId := c.Param("productId")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "productId is required"})
		return
	}

	var results []tables.Result

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
