package db

import "time"

// Product represents a product with its details and relationships.
type Product struct {
	ID            string         `gorm:"type:uuid;primaryKey" json:"id"`
	FullName      string         `json:"fullName"`
	ShortName     string         `json:"shortName"`
	ContactEmail  string         `json:"contactEmail"`
	Relationships []Relationship `gorm:"foreignKey:ObjectIDs;references:ID" json:"relationships"`
}

// Relationship represents a relationship between two objects.
type Relationship struct {
	RelationID       string   `gorm:"type:uuid;primaryKey" json:"relationID"`
	ObjectIDs        []string `gorm:"type:text[]" json:"objectIDs"` // List of two IDs
	RelationshipType string   `json:"relationshipType"`             // e.g., "integration", "dependency", etc.
}

// Integration represents an integration between two products.
type Integration struct {
	ID         string  `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID1 string  `gorm:"type:uuid" json:"productID1"`
	ProductID2 string  `gorm:"type:uuid" json:"productID2"`
	Product1   Product `gorm:"foreignKey:ProductID1"`
	Product2   Product `gorm:"foreignKey:ProductID2"`
}

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

// ResultsRule represents a rule applied to test results.
type ResultsRule struct {
	ID         string       `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID  string       `gorm:"type:uuid" json:"productID"`
	Expression string       `json:"expression"`
	AppliesTo  []string     `gorm:"type:text[]" json:"appliesTo"` // List of types: TestSuite, TestCase, Property, All
	RelationID string       `gorm:"type:uuid" json:"relationID"`
	CreatedAt  time.Time    `json:"createdAt"`
	UpdatedAt  time.Time    `json:"updatedAt"`
	Product    Product      `gorm:"foreignKey:ProductID"`
	Relation   Relationship `gorm:"foreignKey:RelationID"`
}
