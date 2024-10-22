package results

import (
	"hypha/api/internal/db/tables"
	"strings"

	"github.com/gin-gonic/gin"
	"github.com/go-orm/gorm"
	"github.com/lib/pq"
	"github.com/rs/zerolog/log"
)

// FetchResultsByRules fetches test suites and test cases from the database based on the provided rules.
// It filters the fetched data according to the specified expressions in the rules and returns the results.
//
// Parameters:
// - dbConn: A connection to the database.
// - rules: A slice of pointers to ResultsRule, which contains the rules for fetching and filtering the data.
//
// Returns:
// - A slice of gin.H containing the filtered test suites and their associated test cases.
// - An error if any occurs during the database operations.
func FetchResultsByRules(dbConn *gorm.DB, rules []*tables.ResultsRule) ([]gin.H, error) {
	var results []gin.H

	for _, rule := range rules {
		var testSuites []tables.TestSuite
		var testCases []tables.TestCase

		if contains(rule.AppliesTo, "suites") {
			if err := dbConn.Where("product_id = ANY(?)", pq.Array(rule.Relationship.ObjectIDs)).Find(&testSuites).Error; err != nil {
				log.Error().Err(err).Msg("Failed to fetch test suites")
				return nil, err
			}
		}

		if contains(rule.AppliesTo, "cases") {
			if err := dbConn.Where("product_id = ANY(?)", pq.Array(rule.Relationship.ObjectIDs)).Find(&testCases).Error; err != nil {
				log.Error().Err(err).Msg("Failed to fetch test cases")
				return nil, err
			}
		}

		filteredSuites := filterTestSuites(testSuites, rule.Expression)
		filteredCases := filterTestCases(testCases, rule.Expression)

		suiteMap := make(map[string][]tables.TestCase)
		for _, testCase := range filteredCases {
			suiteMap[testCase.TestSuiteID] = append(suiteMap[testCase.TestSuiteID], testCase)
		}

		for _, suite := range filteredSuites {
			if cases, exists := suiteMap[suite.ID]; exists {
				suite.TestCases = cases
			}
			results = append(results, gin.H{
				"suite": suite,
			})
		}
	}

	return results, nil
}

// contains checks if a slice contains a specific item.
//
// Parameters:
// - slice: A slice of strings to search within.
// - item: The string to search for.
//
// Returns:
// - A boolean indicating whether the item is found in the slice.
func contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
}

// filterTestSuites filters test suites based on the provided expression.
//
// Parameters:
// - suites: A slice of TestSuite to filter.
// - expression: A string expression to filter the test suites by.
//
// Returns:
// - A slice of TestSuite that match the expression.
func filterTestSuites(suites []tables.TestSuite, expression string) []tables.TestSuite {
	var filtered []tables.TestSuite
	for _, suite := range suites {
		if matchesExpression(suite.Name, expression) {
			filtered = append(filtered, suite)
		}
	}
	return filtered
}

// filterTestCases filters test cases based on the provided expression.
//
// Parameters:
// - cases: A slice of TestCase to filter.
// - expression: A string expression to filter the test cases by.
//
// Returns:
// - A slice of TestCase that match the expression.
func filterTestCases(cases []tables.TestCase, expression string) []tables.TestCase {
	var filtered []tables.TestCase
	for _, testCase := range cases {
		if matchesExpression(testCase.Name, expression) {
			filtered = append(filtered, testCase)
		}
	}
	return filtered
}

// matchesExpression checks if a name matches the provided expression.
// If the expression starts with '!', it negates the match.
//
// Parameters:
// - name: The name to check against the expression.
// - expression: The expression to match the name against.
//
// Returns:
// - A boolean indicating whether the name matches the expression.
func matchesExpression(name, expression string) bool {
	if strings.HasPrefix(expression, "!") {
		return !strings.HasPrefix(name, expression[1:])
	}
	return strings.HasPrefix(name, expression)
}
