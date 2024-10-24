package tables

import (
	"github.com/lib/pq"
)

// ObjectInterface defines a contract for objects that have an ID.
// Any struct implementing this interface must provide a GetID method that returns a string.
type ObjectInterface interface {
	GetID() string
}

// Product represents a product with its details and relationships.
type Product struct {
	ID            string         `gorm:"type:uuid;primaryKey" json:"id"`
	FullName      string         `json:"fullName"`
	ShortName     string         `json:"shortName"`
	ContactEmail  string         `json:"contactEmail"`
	Relationships []Relationship `gorm:"foreignKey:ObjectIDs;references:ID" json:"relationships"`
}

// GetID returns the ID of the Product.
// This method is required to satisfy the ObjectInterface.
func (p Product) GetID() string {
	return p.ID
}

// Relationship represents a relationship between two objects.
type Relationship struct {
	ID               string            `gorm:"type:uuid;primaryKey" json:"id"`
	ObjectIDs        pq.StringArray    `gorm:"type:text[]" json:"objectIDs"` // List of two IDs
	RelationshipType string            `json:"relationshipType"`             // e.g., "integration", "dependency", etc.
	Objects          []ObjectInterface `gorm:"-" json:"objects"`
}
