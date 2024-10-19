package results

import (
	"bytes"
	"hypha/api/internal/db"
	"strings"
	"time"
)

// ParseJUnitResults parses JUnit test results and stores them in the database.
//
// This function iterates over the provided JUnit test suites, creates corresponding
// result and test suite models, and saves them to the database. It also creates and
// saves properties and test cases associated with each test suite.
//
// Parameters:
// - testSuites: The JUnitTestSuites containing the test results to be parsed.
// - dbOps: The DatabaseOperations interface for interacting with the database.
// - productId: The ID of the product for which the test results are being parsed.
//
// Returns:
// - error: An error if there is any issue during the parsing or saving of the test results.
func ParseJUnitResults(testSuites JUnitTestSuites, dbOps db.DatabaseOperations, productId string) error {
	for _, suite := range testSuites.TestSuites {
		resultModel, err := createResultModel(productId)
		if err != nil {
			return err
		}
		if err := dbOps.Create(&resultModel); err != nil {
			return err
		}

		testSuiteModel, err := createTestSuiteModel(suite, resultModel.ID)
		if err != nil {
			return err
		}
		if err := dbOps.Create(&testSuiteModel); err != nil {
			return err
		}

		if err := createAndSaveProperties(suite.Properties, testSuiteModel.ID, dbOps); err != nil {
			return err
		}

		if err := createAndSaveTestCases(suite.TestCases, testSuiteModel.ID, dbOps); err != nil {
			return err
		}
	}
	return nil
}

// ContainsTestsuitesTag checks if the given XML content contains a <testsuites> tag.
//
// Parameters:
// - xmlContent: The XML content to check as a byte slice.
//
// Returns:
// - bool: A boolean indicating whether the <testsuites> tag is present in the XML content.
func ContainsTestsuitesTag(xmlContent []byte) bool {
	return bytes.Contains(xmlContent, []byte("<testsuites"))
}

// WrapInTestsuitesTag wraps the given XML content in a <testsuites> tag.
//
// Parameters:
// - xmlContent: The XML content to wrap as a byte slice.
//
// Returns:
// - []byte: A new byte slice with the XML content wrapped in a <testsuites> tag.
func WrapInTestsuitesTag(xmlContent []byte) []byte {
	return append([]byte("<testsuites>"), append(xmlContent, []byte("</testsuites>")...)...)
}

// trimLeadingWhitespace removes the leading whitespace from each line of the input text.
// It calculates the minimum indentation level across all lines and removes that amount
// of leading whitespace from each line, preserving any additional whitespace. It also
// removes leading and trailing blank lines.
//
// Parameters:
// - text: The input text with potential leading whitespace.
//
// Returns:
//   - string: A string with the leading whitespace removed from each line based on the minimum
//     indentation level, and with leading and trailing blank lines removed.
func trimLeadingWhitespace(text string) string {
	lines := strings.Split(text, "\n")
	if len(lines) == 0 {
		return text
	}

	// Remove leading and trailing blank lines
	start := 0
	for start < len(lines) && strings.TrimSpace(lines[start]) == "" {
		start++
	}
	end := len(lines) - 1
	for end >= 0 && strings.TrimSpace(lines[end]) == "" {
		end--
	}
	if start > end {
		return ""
	}
	lines = lines[start : end+1]

	minIndent := -1
	for _, line := range lines {
		trimmed := strings.TrimLeft(line, " \t")
		if len(trimmed) == 0 {
			continue
		}
		indent := len(line) - len(trimmed)
		if minIndent == -1 || indent < minIndent {
			minIndent = indent
		}
	}

	for i, line := range lines {
		if len(line) > minIndent {
			lines[i] = line[minIndent:]
		}
	}

	return strings.Join(lines, "\n")
}

// createResultModel creates a new Result model with the given productId.
// It generates a unique ID and sets the current UTC time as the DateReported.
//
// Parameters:
// - productId: The ID of the product for which the result is being created.
//
// Returns:
// - db.Result: The created Result model.
// - error: An error if there is any issue during the creation of the model.
func createResultModel(productId string) (db.Result, error) {
	return db.Result{
		ID:           db.GenerateUniqueID(),
		ProductID:    productId,
		DateReported: time.Now().UTC(),
	}, nil
}

// createTestSuiteModel creates a new TestSuite model from the given JUnitTestSuite and resultID.
// It generates a unique ID for the TestSuite and populates the fields based on the provided suite.
//
// Parameters:
// - suite: The JUnitTestSuite containing the data for the TestSuite model.
// - resultID: The ID of the associated result.
//
// Returns:
// - db.TestSuite: The created TestSuite model.
// - error: An error if there is any issue during the creation of the model.
func createTestSuiteModel(suite JUnitTestSuite, resultID string) (db.TestSuite, error) {
	return db.TestSuite{
		ID:         db.GenerateUniqueID(),
		ResultID:   resultID,
		Name:       suite.Name,
		Tests:      suite.Tests,
		Failures:   suite.Failures,
		Errors:     suite.Errors,
		Skipped:    suite.Skipped,
		Assertions: suite.Assertions,
		Time:       suite.Time,
		File:       suite.File,
		SystemOut:  suite.SystemOut,
		SystemErr:  suite.SystemErr,
	}, nil
}

// createAndSaveProperties creates and saves Property models for the given properties and testSuiteID.
// It iterates over the provided properties, creates a Property model for each, and saves it to the database.
//
// Parameters:
// - properties: A slice of Property structs containing the data for each property.
// - testSuiteID: The ID of the associated test suite.
// - dbOps: The DatabaseOperations interface for interacting with the database.
//
// Returns:
// - error: An error if there is any issue during the creation or saving of the properties.
func createAndSaveProperties(properties []Property, testSuiteID string, dbOps db.DatabaseOperations) error {
	for _, property := range properties {
		value := property.Value
		if value == "" {
			value = trimLeadingWhitespace(property.Text)
		}
		propertyModel := db.Property{
			ID:          db.GenerateUniqueID(),
			TestSuiteID: &testSuiteID,
			Name:        property.Name,
			Value:       value,
		}
		if err := dbOps.Create(&propertyModel); err != nil {
			return err
		}
	}
	return nil
}

// createAndSaveTestCases creates and saves TestCase models for the given test cases and testSuiteID.
// It iterates over the provided test cases, creates a TestCase model for each, and saves it to the database.
// Additionally, it creates and saves properties for each test case.
//
// Parameters:
// - testCases: A slice of JUnitTestCase structs containing the data for each test case.
// - testSuiteID: The ID of the associated test suite.
// - dbOps: The DatabaseOperations interface for interacting with the database.
//
// Returns:
// - error: An error if there is any issue during the creation or saving of the test cases or their properties.
func createAndSaveTestCases(testCases []JUnitTestCase, testSuiteID string, dbOps db.DatabaseOperations) error {
	for _, testCase := range testCases {
		testCaseModel, err := createTestCaseModel(testCase, testSuiteID)
		if err != nil {
			return err
		}
		if err := dbOps.Create(&testCaseModel); err != nil {
			return err
		}

		if err := createAndSaveTestCaseProperties(testCase.Properties, testCaseModel.ID, dbOps); err != nil {
			return err
		}
	}
	return nil
}

// createTestCaseModel creates a new TestCase model from the given JUnitTestCase and testSuiteID.
// It determines the status, message, and type of the test case, generates a unique ID, and populates the fields.
//
// Parameters:
// - testCase: The JUnitTestCase containing the data for the TestCase model.
// - testSuiteID: The ID of the associated test suite.
//
// Returns:
// - db.TestCase: The created TestCase model.
// - error: An error if there is any issue during the creation of the model.
func createTestCaseModel(testCase JUnitTestCase, testSuiteID string) (db.TestCase, error) {
	status, message, testCaseType := determineTestCaseStatus(testCase)

	return db.TestCase{
		ID:          db.GenerateUniqueID(),
		TestSuiteID: testSuiteID,
		ClassName:   testCase.ClassName,
		Name:        testCase.Name,
		Time:        testCase.Time,
		Status:      status,
		Message:     message,
		Type:        testCaseType,
		Assertions:  testCase.Assertions,
		File:        testCase.File,
		Line:        testCase.Line,
		SystemOut:   testCase.SystemOut,
		SystemErr:   testCase.SystemErr,
	}, nil
}

// determineTestCaseStatus determines the status, message, and type of a given JUnitTestCase.
// It checks if the test case has a failure, error, or is skipped, and sets the status accordingly.
//
// Parameters:
// - testCase: The JUnitTestCase for which the status, message, and type need to be determined.
//
// Returns:
// - string: The status of the test case ("pass", "fail", "error", or "skipped").
// - *string: The message associated with the test case status, if any.
// - *string: The type of the test case status, if any.
func determineTestCaseStatus(testCase JUnitTestCase) (string, *string, *string) {
	status := "pass"
	var message *string
	var testCaseType *string

	if testCase.Failure != nil {
		status = "fail"
		message = &testCase.Failure.Message
		testCaseType = &testCase.Failure.Type
	} else if testCase.Error != nil {
		status = "error"
		message = &testCase.Error.Message
		testCaseType = &testCase.Error.Type
	} else if testCase.Skipped != nil {
		status = "skipped"
		message = &testCase.Skipped.Message
		testCaseType = nil
	}

	return status, message, testCaseType
}

// createAndSaveTestCaseProperties creates and saves Property models for the given properties and testCaseID.
// It iterates over the provided properties, creates a Property model for each, and saves it to the database.
//
// Parameters:
// - properties: A slice of Property structs containing the data for each property.
// - testCaseID: The ID of the associated test case.
// - dbOps: The DatabaseOperations interface for interacting with the database.
//
// Returns:
// - error: An error if there is any issue during the creation or saving of the properties.
func createAndSaveTestCaseProperties(properties []Property, testCaseID string, dbOps db.DatabaseOperations) error {
	for _, property := range properties {
		value := property.Value
		if value == "" {
			value = trimLeadingWhitespace(property.Text)
		}
		propertyModel := db.Property{
			ID:         db.GenerateUniqueID(),
			TestCaseID: &testCaseID,
			Name:       property.Name,
			Value:      value,
		}
		if err := dbOps.Create(&propertyModel); err != nil {
			return err
		}
	}
	return nil
}
