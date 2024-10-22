package queries

import (
	"hypha/api/internal/db"
	"hypha/api/internal/db/tables"
)

// FetchRulesByRelationID retrieves results rules by their relation ID from the database.
//
// Parameters:
// - dbOps: The database operations interface used for database interactions.
// - relationID: The relation ID to filter by.
//
// Returns:
// - []*tables.ResultsRule: A slice of results rule objects if found.
// - error: An error if any database operation fails.
func FetchRulesByRelationID(dbOps db.DatabaseOperations, relationID string) ([]*tables.ResultsRule, error) {
	var rules []*tables.ResultsRule
	err := dbOps.Connection().Where("relationship_id = ?", relationID).Find(&rules).Error
	if err != nil {
		return nil, err
	}
	return rules, nil
}
