package report

import (
	"encoding/xml"
	"hypha/api/internal/db/ops"
	"hypha/api/internal/db/tables"
	"io/ioutil"
	"net/http"

	"hypha/api/internal/utils/results/parse"
	"hypha/api/internal/utils/results/structs"

	"github.com/gin-gonic/gin"
)

func InitReportRoutes(router *gin.RouterGroup, dbOperations ops.DatabaseOperations) {
	router.POST("/results", func(c *gin.Context) {
		ReportResults(c, dbOperations)
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

	if err := xml.Unmarshal(byteValue, &junitTestSuites); err != nil {
		c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid XML format. Not JUnit"})
		return
	}

	if err := parse.ParseJUnitResults(junitTestSuites, dbOperations, productId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
