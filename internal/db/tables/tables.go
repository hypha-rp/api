package tables

import "time"

type Product struct {
	ID           string        `gorm:"type:uuid;primaryKey" json:"id"`
	FullName     string        `json:"fullName"`
	ShortName    string        `json:"shortName"`
	ContactEmail string        `json:"contactEmail"`
	Integrations []Integration `gorm:"foreignKey:ProductID1;foreignKey:ProductID2"`
}

type Integration struct {
	ID         string  `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID1 string  `gorm:"type:uuid;" json:"productID1"`
	ProductID2 string  `gorm:"type:uuid" json:"productID2"`
	Product1   Product `gorm:"foreignKey:ProductID1"`
	Product2   Product `gorm:"foreignKey:ProductID2"`
}

type Result struct {
	ID           string      `gorm:"type:uuid;primaryKey" json:"id"`
	ProductID    string      `json:"productID"`
	TestSuites   []TestSuite `gorm:"foreignKey:ResultID"`
	DateReported time.Time   `json:"dateReported"`
}

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

type Property struct {
	ID          string  `gorm:"type:uuid;primaryKey" json:"id"`
	TestSuiteID *string `json:"testSuiteID"`
	TestCaseID  *string `json:"testCaseID"`
	Name        string  `json:"name"`
	Value       string  `json:"value"`
}
