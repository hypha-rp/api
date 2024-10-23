package tables

import (
	"time"

	"github.com/lib/pq"
)

// ResultsRule represents a rule applied to test results.
type ResultsRule struct {
	ID             string         `gorm:"type:uuid;primaryKey" json:"id"`
	Expression     string         `json:"expression"`
	AppliesTo      pq.StringArray `gorm:"type:text[]" json:"appliesTo"` // List of types: suite, case
	RelationshipID string         `gorm:"type:uuid" json:"relationshipId"`
	Relationship   Relationship   `gorm:"foreignKey:RelationshipID"`
	CreatedAt      time.Time      `json:"createdAt"`
	UpdatedAt      time.Time      `json:"updatedAt"`
}
