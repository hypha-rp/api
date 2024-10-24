package queries

import (
	"hypha/api/internal/db/tables"
	"hypha/api/internal/utils"

	"github.com/go-orm/gorm"
	"github.com/lib/pq"
)

func FetchResultsByRules(dbConn *gorm.DB, rules []*tables.ResultsRule) ([]tables.Result, error) {
	var resultsMap = make(map[string]tables.Result)

	for _, rule := range rules {
		var viewResults []tables.TestResultsView

		// Fetch results from the view based on relationship object IDs
		if err := dbConn.Where("product_id = ANY(?)", pq.Array(rule.Relationship.ObjectIDs)).Find(&viewResults).Error; err != nil {
			return nil, err
		}

		// Fetch properties for test suites and test cases
		var suiteProperties []tables.Property
		var caseProperties []tables.Property
		suiteIDs := make([]string, 0)
		caseIDs := make([]string, 0)
		for _, vr := range viewResults {
			if vr.TestSuiteID != "" {
				suiteIDs = append(suiteIDs, vr.TestSuiteID)
			}
			if vr.TestCaseID != "" {
				caseIDs = append(caseIDs, vr.TestCaseID)
			}
		}

		if len(suiteIDs) > 0 {
			if err := dbConn.Where("test_suite_id IN (?)", suiteIDs).Find(&suiteProperties).Error; err != nil {
				return nil, err
			}
		}

		if len(caseIDs) > 0 {
			if err := dbConn.Where("test_case_id IN (?)", caseIDs).Find(&caseProperties).Error; err != nil {
				return nil, err
			}
		}

		// Filter the results based on the rule's expression
		filteredSuites := make(map[string]tables.TestSuite)
		filteredCases := make(map[string][]tables.TestCase)

		for _, vr := range viewResults {
			if utils.Contains(rule.AppliesTo, "suite") && utils.MatchesExpression(vr.TestSuiteName, rule.Expression) {
				if _, exists := filteredSuites[vr.TestSuiteID]; !exists {
					filteredSuites[vr.TestSuiteID] = tables.TestSuite{
						ID:         vr.TestSuiteID,
						ResultID:   vr.ResultID,
						Name:       vr.TestSuiteName,
						Tests:      vr.TestSuiteTests,
						Failures:   vr.TestSuiteFailures,
						Errors:     vr.TestSuiteErrors,
						Skipped:    vr.TestSuiteSkipped,
						Assertions: vr.TestSuiteAssertions,
						Time:       vr.TestSuiteTime,
						File:       vr.TestSuiteFile,
						SystemOut:  vr.TestSuiteSystemOut,
						SystemErr:  vr.TestSuiteSystemErr,
						TestCases:  []tables.TestCase{},
						Properties: []tables.Property{},
					}
				}
				// Ensure all test cases for the matching suite are included
				if vr.TestCaseID != "" {
					filteredCases[vr.TestSuiteID] = append(filteredCases[vr.TestSuiteID], tables.TestCase{
						ID:          vr.TestCaseID,
						TestSuiteID: vr.TestSuiteID,
						ClassName:   vr.TestCaseClassName,
						Name:        vr.TestCaseName,
						Time:        vr.TestCaseTime,
						Status:      vr.TestCaseStatus,
						Message:     vr.TestCaseMessage,
						Type:        vr.TestCaseType,
						Assertions:  vr.TestCaseAssertions,
						File:        vr.TestCaseFile,
						Line:        vr.TestCaseLine,
						SystemOut:   vr.TestCaseSystemOut,
						SystemErr:   vr.TestCaseSystemErr,
						Properties:  []tables.Property{},
					})
				}
			}
			if utils.Contains(rule.AppliesTo, "case") && utils.MatchesExpression(vr.TestCaseName, rule.Expression) {
				filteredCases[vr.TestSuiteID] = append(filteredCases[vr.TestSuiteID], tables.TestCase{
					ID:          vr.TestCaseID,
					TestSuiteID: vr.TestSuiteID,
					ClassName:   vr.TestCaseClassName,
					Name:        vr.TestCaseName,
					Time:        vr.TestCaseTime,
					Status:      vr.TestCaseStatus,
					Message:     vr.TestCaseMessage,
					Type:        vr.TestCaseType,
					Assertions:  vr.TestCaseAssertions,
					File:        vr.TestCaseFile,
					Line:        vr.TestCaseLine,
					SystemOut:   vr.TestCaseSystemOut,
					SystemErr:   vr.TestCaseSystemErr,
					Properties:  []tables.Property{},
				})
			}
		}

		// Ensure all test suites that have test cases are included
		for suiteID, cases := range filteredCases {
			if suite, exists := filteredSuites[suiteID]; exists {
				suite.TestCases = cases
				filteredSuites[suiteID] = suite
			} else {
				// Populate suite data when only test cases are matched
				for _, vr := range viewResults {
					if vr.TestSuiteID == suiteID {
						filteredSuites[suiteID] = tables.TestSuite{
							ID:         vr.TestSuiteID,
							ResultID:   vr.ResultID,
							Name:       vr.TestSuiteName,
							Tests:      vr.TestSuiteTests,
							Failures:   vr.TestSuiteFailures,
							Errors:     vr.TestSuiteErrors,
							Skipped:    vr.TestSuiteSkipped,
							Assertions: vr.TestSuiteAssertions,
							Time:       vr.TestSuiteTime,
							File:       vr.TestSuiteFile,
							SystemOut:  vr.TestSuiteSystemOut,
							SystemErr:  vr.TestSuiteSystemErr,
							TestCases:  cases,
							Properties: []tables.Property{},
						}
						break
					}
				}
			}
		}

		// Add properties to the respective test suites and test cases
		for _, prop := range suiteProperties {
			if suite, exists := filteredSuites[*prop.TestSuiteID]; exists {
				suite.Properties = append(suite.Properties, prop)
				filteredSuites[*prop.TestSuiteID] = suite
			}
		}

		for _, prop := range caseProperties {
			for suiteID, cases := range filteredCases {
				for i, testCase := range cases {
					if testCase.ID == *prop.TestCaseID {
						testCase.Properties = append(testCase.Properties, prop)
						filteredCases[suiteID][i] = testCase
					}
				}
			}
		}

		// Combine filtered suites into their respective results
		for _, suite := range filteredSuites {
			if result, exists := resultsMap[suite.ResultID]; exists {
				result.TestSuites = append(result.TestSuites, suite)
				resultsMap[suite.ResultID] = result
			} else {
				// Fetch the actual result to get the accurate DateReported
				var actualResult tables.Result
				if err := dbConn.First(&actualResult, "id = ?", suite.ResultID).Error; err != nil {
					return nil, err
				}

				resultsMap[suite.ResultID] = tables.Result{
					ID:           suite.ResultID,
					ProductID:    viewResults[0].ProductID,
					TestSuites:   []tables.TestSuite{suite},
					DateReported: actualResult.DateReported,
				}
			}
		}
	}

	// Convert resultsMap to a slice
	var results []tables.Result
	for _, result := range resultsMap {
		results = append(results, result)
	}

	return results, nil
}
