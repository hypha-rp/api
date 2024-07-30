package report

import (
	"encoding/xml"
	"hypha/api/internal/db/ops"
	"hypha/api/internal/db/tables"
	"io/ioutil"
	"net/http"
	"strconv"

	"github.com/gin-gonic/gin"
)

type Assemblies struct {
	XMLName    xml.Name   `xml:"assemblies"`
	Assemblies []Assembly `xml:"assembly"`
}

type Assembly struct {
	XMLName       xml.Name     `xml:"assembly"`
	ID            string       `xml:"id,attr"`
	Name          string       `xml:"name,attr"`
	TestFramework string       `xml:"test-framework,attr"`
	RunDate       string       `xml:"run-date,attr"`
	RunTime       string       `xml:"run-time,attr"`
	Total         int          `xml:"total,attr"`
	Passed        int          `xml:"passed,attr"`
	Failed        int          `xml:"failed,attr"`
	Skipped       int          `xml:"skipped,attr"`
	Time          float64      `xml:"time,attr"`
	Collections   []Collection `xml:"collection"`
}

type Collection struct {
	XMLName xml.Name `xml:"collection"`
	ID      string   `xml:"id,attr"`
	Total   int      `xml:"total,attr"`
	Passed  int      `xml:"passed,attr"`
	Failed  int      `xml:"failed,attr"`
	Skipped int      `xml:"skipped,attr"`
	Name    string   `xml:"name,attr"`
	Tests   []Test   `xml:"test"`
}

type Test struct {
	XMLName xml.Name `xml:"test"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Method  string   `xml:"method,attr"`
	Time    float64  `xml:"time,attr"`
	Result  string   `xml:"result,attr"`
	Traits  []Trait  `xml:"traits>trait"`
}

type Trait struct {
	XMLName xml.Name `xml:"trait"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

func InitReportRoutes(router *gin.RouterGroup, dbOperations ops.DatabaseOperations) {
	router.POST("/results", func(c *gin.Context) {
		productId := c.PostForm("productId")
		var product tables.Product
		var assemblies Assemblies
		var err error

		if productId == "" {
			c.JSON(http.StatusBadRequest, gin.H{"error": "productId is required"})
			return
		}

		if dbOperations == nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": "Database operations not initialized"})
			return
		}

		if err = dbOperations.First(&product, "id = ?", productId); err != nil {
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
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid XML"})
			return
		}

		parsedProductId, err := strconv.ParseUint(productId, 10, 32)
		if err != nil {
			c.JSON(http.StatusBadRequest, gin.H{"error": "Invalid productId"})
			return
		}

		if err := ParseResults(assemblies, dbOperations, uint(parsedProductId)); err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}

		c.JSON(http.StatusOK, gin.H{"status": "success"})
	})
}

func ParseResults(assemblies Assemblies, dbOperations ops.DatabaseOperations, productId uint) error {
	for _, assembly := range assemblies.Assemblies {
		assemblyModel := tables.Assembly{
			ID:            assembly.ID,
			Name:          assembly.Name,
			TestFramework: assembly.TestFramework,
			RunDate:       assembly.RunDate,
			RunTime:       assembly.RunTime,
			Total:         assembly.Total,
			Passed:        assembly.Passed,
			Failed:        assembly.Failed,
			Skipped:       assembly.Skipped,
			Time:          assembly.Time,
			ProductID:     productId,
		}
		if err := dbOperations.Create(&assemblyModel); err != nil {
			return err
		}

		for _, collection := range assembly.Collections {
			collectionModel := tables.Collection{
				ID:         collection.ID,
				AssemblyID: assemblyModel.ID,
				Total:      collection.Total,
				Passed:     collection.Passed,
				Failed:     collection.Failed,
				Skipped:    collection.Skipped,
				Name:       collection.Name,
			}
			if err := dbOperations.Create(&collectionModel); err != nil {
				return err
			}

			for _, test := range collection.Tests {
				testModel := tables.Test{
					ID:           test.ID,
					CollectionID: collectionModel.ID,
					Name:         test.Name,
					Type:         test.Type,
					Method:       test.Method,
					Time:         test.Time,
					Result:       test.Result,
				}
				if err := dbOperations.Create(&testModel); err != nil {
					return err
				}

				for _, trait := range test.Traits {
					traitModel := tables.Trait{
						TestID: testModel.ID,
						Name:   trait.Name,
						Value:  trait.Value,
					}
					if err := dbOperations.Create(&traitModel); err != nil {
						return err
					}
				}
			}
		}
	}
	return nil
}
