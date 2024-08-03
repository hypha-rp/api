package report

import (
	"encoding/xml"
	"hypha/api/internal/db/ops"
	"hypha/api/internal/db/tables"
	"io/ioutil"
	"net/http"

	"hypha/api/internal/utils/results/conversion"
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
	var assemblies structs.Assemblies
	var product tables.Product

	if productId := c.PostForm("productId"); productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "productId is required"})
		return
	}

	if dbOperations == nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": "Database operations not initialized"})
		return
	}

	productId := c.PostForm("productId")
	if productId == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "productId is required"})
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

	if err := xml.Unmarshal(byteValue, &assemblies); err != nil {
		var junitTestSuites structs.JUnitTestSuites
		if err := xml.Unmarshal(byteValue, &junitTestSuites); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid XML format. Not XUnit or JUnit"})
			return
		}

		xunitFile, err := conversion.ConvertJUnitToXUnit(file.Filename)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to convert JUnit to XUnit"})
			return
		}

		byteValue, err = ioutil.ReadFile(xunitFile)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Failed to read converted XUnit file"})
			return
		}

		if err := xml.Unmarshal(byteValue, &assemblies); err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid XML format after conversion"})
			return
		}
	}

	if err := parse.ParseXUnitResults(assemblies, dbOperations, productId); err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}

	c.JSON(http.StatusOK, gin.H{"status": "success"})
}
