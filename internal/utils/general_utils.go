package utils

import "strings"

// contains checks if a slice contains a specific item.
//
// Parameters:
// - slice: A slice of strings to search within.
// - item: The string to search for.
//
// Returns:
// - A boolean indicating whether the item is found in the slice.
func Contains(slice []string, item string) bool {
	for _, v := range slice {
		if v == item {
			return true
		}
	}
	return false
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
func MatchesExpression(value, expression string) bool {
	if strings.HasPrefix(expression, "!") {
		return !strings.HasPrefix(value, expression[1:])
	}
	return strings.HasPrefix(value, expression)
}
