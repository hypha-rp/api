package tables

import (
	"time"
)

// Result represents a test result for a product.
type Result struct {
	ID           string      `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID    string      `json:"productID"`
	TestSuites   []TestSuite `gorm:"foreignKey:ResultID"`
	DateReported time.Time   `json:"dateReported"`
}

// TestSuite represents a suite of tests within a test result.
type TestSuite struct {
	ID         string     `gorm:"type:uuid;primaryKey" json:"id"`
	ResultID   string     `json:"resultID"`
	Name       string     `json:"name"`
	Tests      int        `json:"tests"`
	Failures   int        `json:"failures"`
	Errors     int        `json:"errors"`
	Skipped    int        `json:"skipped"`
	Assertions int        `json:"assertions"`
	Time       float64    `json:"time"`
	File       string     `json:"file"`
	TestCases  []TestCase `gorm:"foreignKey:TestSuiteID"`
	Properties []Property `gorm:"foreignKey:TestSuiteID"`
	SystemOut  string     `json:"systemOut"`
	SystemErr  string     `json:"systemErr"`
}

// TestCase represents an individual test case within a test suite.
type TestCase struct {
	ID          string     `gorm:"type:uuid;primaryKey" json:"id"`
	TestSuiteID string     `json:"testSuiteID"`
	ClassName   string     `json:"className"`
	Name        string     `json:"name"`
	Time        float64    `json:"time"`
	Status      string     `json:"status"`
	Message     *string    `json:"message"`
	Type        *string    `json:"type"`
	Assertions  int        `json:"assertions"`
	File        string     `json:"file"`
	Line        int        `json:"line"`
	Properties  []Property `gorm:"foreignKey:TestCaseID"`
	SystemOut   string     `json:"systemOut"`
	SystemErr   string     `json:"systemErr"`
}

// Property represents a property associated with a test suite or test case.
type Property struct {
	ID          string  `gorm:"type:uuid;primaryKey" json:"id"`
	TestSuiteID *string `json:"testSuiteID"`
	TestCaseID  *string `json:"testCaseID"`
	Name        string  `json:"name"`
	Value       string  `json:"value"`
}

// TestResultsView represents the view for test results.
type TestResultsView struct {
	ResultID            string  `json:"result_id"`
	ProductID           string  `json:"product_id"`
	TestSuiteID         string  `json:"test_suite_id"`
	TestSuiteName       string  `json:"test_suite_name"`
	TestSuiteTests      int     `json:"test_suite_tests"`
	TestSuiteFailures   int     `json:"test_suite_failures"`
	TestSuiteErrors     int     `json:"test_suite_errors"`
	TestSuiteSkipped    int     `json:"test_suite_skipped"`
	TestSuiteAssertions int     `json:"test_suite_assertions"`
	TestSuiteTime       float64 `json:"test_suite_time"`
	TestSuiteFile       string  `json:"test_suite_file"`
	TestSuiteSystemOut  string  `json:"test_suite_system_out"`
	TestSuiteSystemErr  string  `json:"test_suite_system_err"`
	TestCaseID          string  `json:"test_case_id"`
	TestCaseName        string  `json:"test_case_name"`
	TestCaseClassName   string  `json:"test_case_class_name"`
	TestCaseTime        float64 `json:"test_case_time"`
	TestCaseStatus      string  `json:"test_case_status"`
	TestCaseMessage     *string `json:"test_case_message"`
	TestCaseType        *string `json:"test_case_type"`
	TestCaseAssertions  int     `json:"test_case_assertions"`
	TestCaseFile        string  `json:"test_case_file"`
	TestCaseLine        int     `json:"test_case_line"`
	TestCaseSystemOut   string  `json:"test_case_system_out"`
	TestCaseSystemErr   string  `json:"test_case_system_err"`
}

// TableName sets the insert table name for this struct type
func (TestResultsView) TableName() string {
	return "test_results_view"
}
