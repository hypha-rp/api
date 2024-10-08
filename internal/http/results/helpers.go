package results

import (
	"hypha/api/internal/db/tables"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-orm/gorm"
)

func logErrorAndRespond(context *gin.Context, message string, err error) {
	log.Error().Msgf("%s: %v", message, err)
	context.JSON(http.StatusInternalServerError, gin.H{"error": "There was a problem processing your request"})
}

func createResultMap(testSuites []tables.TestSuite) map[string][]tables.TestSuite {
	resultMap := make(map[string][]tables.TestSuite)
	for _, testSuite := range testSuites {
		resultID := testSuite.ResultID
		resultMap[resultID] = append(resultMap[resultID], testSuite)
	}
	return resultMap
}

func fetchResultsAndProducts(db *gorm.DB, resultMap map[string][]tables.TestSuite) ([]gin.H, error) {
	var results []gin.H
	for resultID, testSuites := range resultMap {
		var result tables.Result
		err := db.Where("id = ?", resultID).First(&result).Error
		if err != nil {
			return nil, err
		}
		result.TestSuites = testSuites

		var product tables.Product
		err = db.Where("id = ?", result.ProductID).First(&product).Error
		if err != nil {
			return nil, err
		}

		results = append(results, gin.H{
			"id":           result.ID,
			"productID":    result.ProductID,
			"productName":  product.FullName,
			"TestSuites":   result.TestSuites,
			"dateReported": result.DateReported,
		})
	}
	return results, nil
}

func getTestSuiteAndCaseIDs(db *gorm.DB, integrationID string) ([]string, []string, error) {
	var testSuiteIDs []string
	var testCaseIDs []string

	err := db.Table("properties").
		Where("properties.name = ? AND properties.value::text = ? AND properties.test_suite_id IS NOT NULL", "hypha.integration", integrationID).
		Pluck("test_suite_id::text", &testSuiteIDs).Error

	if err != nil {
		log.Error().Msgf("Database query error in getTestSuiteAndCaseIDs (test_suite_ids): %v", err)
		return nil, nil, err
	}

	err = db.Table("properties").
		Where("properties.name = ? AND properties.value::text = ? AND properties.test_case_id IS NOT NULL", "hypha.integration", integrationID).
		Pluck("test_case_id::text", &testCaseIDs).Error

	if err != nil {
		log.Error().Msgf("Database query error in getTestSuiteAndCaseIDs (test_case_ids): %v", err)
		return nil, nil, err
	}

	return testSuiteIDs, testCaseIDs, nil
}

func getTestSuites(db *gorm.DB, testSuiteIDs, testCaseIDs []string) ([]tables.TestSuite, error) {
	var testSuites []tables.TestSuite

	err := db.Where("id::text IN (?) OR id::text IN (SELECT test_suite_id::text FROM test_cases WHERE id::text IN (?))", testSuiteIDs, testCaseIDs).
		Preload("TestCases").
		Preload("TestCases.Properties").
		Preload("Properties").
		Find(&testSuites).Error

	if err != nil {
		log.Error().Msgf("Database query error in getTestSuites: %v", err)
		return nil, err
	}

	return testSuites, nil
}

func filterTestCases(testSuites []tables.TestSuite, integrationID string) {
	for i := range testSuites {
		integrationAtSuiteLevel := false
		for _, property := range testSuites[i].Properties {
			if property.Name == "hypha.integration" && property.Value == integrationID {
				integrationAtSuiteLevel = true
				break
			}
		}

		if integrationAtSuiteLevel {
			continue
		}

		filteredTestCases := []tables.TestCase{}
		for _, testCase := range testSuites[i].TestCases {
			for _, property := range testCase.Properties {
				if property.Name == "hypha.integration" && property.Value == integrationID {
					filteredTestCases = append(filteredTestCases, testCase)
					break
				}
			}
		}
		testSuites[i].TestCases = filteredTestCases
	}
}
