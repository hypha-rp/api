package utils

import "strings"

// Contains checks if a slice contains a specific item.
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

// MatchesExpression checks if a value matches the provided expression.
// If the expression starts with '!', it negates the match.
// Supports wildcard '*' for matching any sequence of characters.
//
// Parameters:
// - value: The value to check against the expression.
// - expression: The expression to match the value against.
//
// Returns:
// - A boolean indicating whether the value matches the expression.
func MatchesExpression(value, expression string) bool {
	if strings.HasPrefix(expression, "!") {
		return !MatchesWildcard(value, expression[1:])
	}
	return MatchesWildcard(value, expression)
}

// matchesWildcard checks if a value matches a wildcard pattern.
// Supports '*' as a wildcard character.
//
// Parameters:
// - value: The value to check against the pattern.
// - pattern: The wildcard pattern to match the value against.
//
// Returns:
// - A boolean indicating whether the value matches the pattern.
func MatchesWildcard(value, pattern string) bool {
	parts := strings.Split(pattern, "*")
	for _, part := range parts {
		idx := strings.Index(value, part)
		if idx == -1 {
			return false
		}
		value = value[idx+len(part):]
	}
	return true
}
