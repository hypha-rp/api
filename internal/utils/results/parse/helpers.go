package parse

import (
	"hypha/api/internal/db/ops"
	"hypha/api/internal/db/tables"
	"hypha/api/internal/utils/results/structs"
	"strings"
	"time"
)

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
// - tables.Result: The created Result model.
// - error: An error if there is any issue during the creation of the model.
func createResultModel(productId string) (tables.Result, error) {
	return tables.Result{
		ID:           ops.GenerateUniqueID(),
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
// - tables.TestSuite: The created TestSuite model.
// - error: An error if there is any issue during the creation of the model.
func createTestSuiteModel(suite structs.JUnitTestSuite, resultID string) (tables.TestSuite, error) {
	return tables.TestSuite{
		ID:         ops.GenerateUniqueID(),
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
// - dbOperations: The DatabaseOperations interface for interacting with the database.
//
// Returns:
// - error: An error if there is any issue during the creation or saving of the properties.
func createAndSaveProperties(properties []structs.Property, testSuiteID string, dbOperations ops.DatabaseOperations) error {
	for _, property := range properties {
		value := property.Value
		if value == "" {
			value = trimLeadingWhitespace(property.Text)
		}
		propertyModel := tables.Property{
			ID:          ops.GenerateUniqueID(),
			TestSuiteID: &testSuiteID,
			Name:        property.Name,
			Value:       value,
		}
		if err := dbOperations.Create(&propertyModel); err != nil {
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
// - dbOperations: The DatabaseOperations interface for interacting with the database.
//
// Returns:
// - error: An error if there is any issue during the creation or saving of the test cases or their properties.
func createAndSaveTestCases(testCases []structs.JUnitTestCase, testSuiteID string, dbOperations ops.DatabaseOperations) error {
	for _, testCase := range testCases {
		testCaseModel, err := createTestCaseModel(testCase, testSuiteID)
		if err != nil {
			return err
		}
		if err := dbOperations.Create(&testCaseModel); err != nil {
			return err
		}

		if err := createAndSaveTestCaseProperties(testCase.Properties, testCaseModel.ID, dbOperations); err != nil {
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
// - tables.TestCase: The created TestCase model.
// - error: An error if there is any issue during the creation of the model.
func createTestCaseModel(testCase structs.JUnitTestCase, testSuiteID string) (tables.TestCase, error) {
	status, message, testCaseType := determineTestCaseStatus(testCase)

	return tables.TestCase{
		ID:          ops.GenerateUniqueID(),
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
func determineTestCaseStatus(testCase structs.JUnitTestCase) (string, *string, *string) {
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
// - dbOperations: The DatabaseOperations interface for interacting with the database.
//
// Returns:
// - error: An error if there is any issue during the creation or saving of the properties.
func createAndSaveTestCaseProperties(properties []structs.Property, testCaseID string, dbOperations ops.DatabaseOperations) error {
	for _, property := range properties {
		value := property.Value
		if value == "" {
			value = trimLeadingWhitespace(property.Text)
		}
		propertyModel := tables.Property{
			ID:         ops.GenerateUniqueID(),
			TestCaseID: &testCaseID,
			Name:       property.Name,
			Value:      value,
		}
		if err := dbOperations.Create(&propertyModel); err != nil {
			return err
		}
	}
	return nil
}
