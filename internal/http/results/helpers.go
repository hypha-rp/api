package results

import (
	"hypha/api/internal/db"
	"net/http"

	"github.com/gin-gonic/gin"
	"github.com/go-orm/gorm"
)

// logErrorAndRespond logs an error message and sends a JSON response with an internal server error status.
// Parameters:
// - context: The Gin context to use for sending the response.
// - message: The error message to log.
// - err: The error to log.
func logErrorAndRespond(context *gin.Context, message string, err error) {
	log.Error().Msgf("%s: %v", message, err)
	context.JSON(http.StatusInternalServerError, gin.H{"error": "There was a problem processing your request"})
}

// createResultMap creates a map of result IDs to their corresponding test suites.
// Parameters:
// - testSuites: A slice of TestSuite objects.
// Returns:
// - A map where the keys are result IDs and the values are slices of TestSuite objects.
func createResultMap(testSuites []db.TestSuite) map[string][]db.TestSuite {
	resultMap := make(map[string][]db.TestSuite)
	for _, testSuite := range testSuites {
		resultID := testSuite.ResultID
		resultMap[resultID] = append(resultMap[resultID], testSuite)
	}
	return resultMap
}

// fetchResultsAndProducts fetches results and their associated products from the database.
// Parameters:
// - db: The GORM database connection.
// - resultMap: A map of result IDs to their corresponding test suites.
// Returns:
// - A slice of Gin H maps containing result and product information.
// - An error if any database operation fails.
func fetchResultsAndProducts(dbConn *gorm.DB, resultMap map[string][]db.TestSuite) ([]gin.H, error) {
	var results []gin.H
	for resultID, testSuites := range resultMap {
		var result db.Result
		err := dbConn.Where("id = ?", resultID).First(&result).Error
		if err != nil {
			return nil, err
		}
		result.TestSuites = testSuites

		var product db.Product
		err = dbConn.Where("id = ?", result.ProductID).First(&product).Error
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

// getTestSuiteAndCaseIDs retrieves test suite and test case IDs associated with a given integration ID.
// Parameters:
// - db: The GORM database connection.
// - integrationID: The integration ID to filter by.
// Returns:
// - A slice of test suite IDs.
// - A slice of test case IDs.
// - An error if any database operation fails.
func getTestSuiteAndCaseIDs(db *gorm.DB, integrationID string) ([]string, []string, error) {
	var testSuiteIDs []string
	var testCaseIDs []string

	// Retrieve test suite IDs associated with the integration ID
	err := db.Table("properties").
		Where("properties.name = ? AND properties.value::text = ? AND properties.test_suite_id IS NOT NULL", "hypha.integration", integrationID).
		Pluck("test_suite_id::text", &testSuiteIDs).Error

	if err != nil {
		log.Error().Msgf("Database query error in getTestSuiteAndCaseIDs (test_suite_ids): %v", err)
		return nil, nil, err
	}

	// Retrieve test case IDs associated with the integration ID
	err = db.Table("properties").
		Where("properties.name = ? AND properties.value::text = ? AND properties.test_case_id IS NOT NULL", "hypha.integration", integrationID).
		Pluck("test_case_id::text", &testCaseIDs).Error

	if err != nil {
		log.Error().Msgf("Database query error in getTestSuiteAndCaseIDs (test_case_ids): %v", err)
		return nil, nil, err
	}

	// If integration is at the suite level, retrieve test case IDs for those suites
	if len(testSuiteIDs) > 0 {
		err = db.Table("test_cases").
			Where("test_suite_id IN (?)", testSuiteIDs).
			Pluck("id::text", &testCaseIDs).Error

		if err != nil {
			log.Error().Msgf("Database query error in getTestSuiteAndCaseIDs (test_cases for suites): %v", err)
			return nil, nil, err
		}
	}

	// Retrieve test case IDs associated with the test suites where the hypha.integration property matches the integration ID
	if len(testSuiteIDs) > 0 {
		err = db.Table("properties").
			Where("properties.name = ? AND properties.value::text = ? AND properties.test_suite_id IN (?) AND properties.test_case_id IS NOT NULL", "hypha.integration", integrationID, testSuiteIDs).
			Pluck("test_case_id::text", &testCaseIDs).Error

		if err != nil {
			log.Error().Msgf("Database query error in getTestSuiteAndCaseIDs (test_case_ids for suites): %v", err)
			return nil, nil, err
		}
	}

	return testSuiteIDs, testCaseIDs, nil
}

// getTestSuites retrieves test suites and their associated test cases and properties from the database.
// Parameters:
// - db: The GORM database connection.
// - testSuiteIDs: A slice of test suite IDs to filter by.
// - testCaseIDs: A slice of test case IDs to filter by.
// Returns:
// - A slice of TestSuite objects.
// - An error if any database operation fails.
func getTestSuites(dbConn *gorm.DB, testSuiteIDs, testCaseIDs []string) ([]db.TestSuite, error) {
	var testSuites []db.TestSuite

	err := dbConn.Where("id::text IN (?) OR id::text IN (SELECT test_suite_id::text FROM test_cases WHERE id::text IN (?))", testSuiteIDs, testCaseIDs).
		Preload("TestCases").
		Preload("TestCases.Properties").
		Preload("Properties").
		Find(&testSuites).Error

	if err != nil {
		return nil, err
	}

	return testSuites, nil
}

// filterTestCases filters test cases within test suites based on the integration ID.
// Parameters:
// - testSuites: A slice of TestSuite objects to filter.
// - integrationID: The integration ID to filter by.
func filterTestCases(testSuites []db.TestSuite, integrationID string) {
	for i := range testSuites {
		var filteredTestCases []db.TestCase
		suiteHasIntegration := false

		for _, property := range testSuites[i].Properties {
			if property.Name == "hypha.integration" && property.Value == integrationID {
				suiteHasIntegration = true
				break
			}
		}

		for _, testCase := range testSuites[i].TestCases {
			caseHasIntegration := false
			for _, property := range testCase.Properties {
				if property.Name == "hypha.integration" && property.Value == integrationID {
					caseHasIntegration = true
					break
				}
			}
			if caseHasIntegration || suiteHasIntegration {
				filteredTestCases = append(filteredTestCases, testCase)
			}
		}
		testSuites[i].TestCases = filteredTestCases
	}
}
