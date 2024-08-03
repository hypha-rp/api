package structs

import (
	"encoding/xml"
)

type Assemblies struct {
	XMLName    xml.Name   `xml:"assemblies"`
	Assemblies []Assembly `xml:"assembly"`
}

type Assembly struct {
	XMLName       xml.Name     `xml:"assembly"`
	ID            string       `xml:"id,attr"`
	Name          string       `xml:"name,attr"`
	TestFramework string       `xml:"test-framework,attr"`
	RunDate       string       `xml:"run-date,attr"`
	RunTime       string       `xml:"run-time,attr"`
	Total         int          `xml:"total,attr"`
	Passed        int          `xml:"passed,attr"`
	Failed        int          `xml:"failed,attr"`
	Skipped       int          `xml:"skipped,attr"`
	Time          float64      `xml:"time,attr"`
	Collections   []Collection `xml:"collection"`
}

type Collection struct {
	XMLName xml.Name `xml:"collection"`
	ID      string   `xml:"id,attr"`
	Total   int      `xml:"total,attr"`
	Passed  int      `xml:"passed,attr"`
	Failed  int      `xml:"failed,attr"`
	Skipped int      `xml:"skipped,attr"`
	Name    string   `xml:"name,attr"`
	Tests   []Test   `xml:"test"`
}

type Test struct {
	XMLName xml.Name `xml:"test"`
	ID      string   `xml:"id,attr"`
	Name    string   `xml:"name,attr"`
	Type    string   `xml:"type,attr"`
	Method  string   `xml:"method,attr"`
	Time    float64  `xml:"time,attr"`
	Result  string   `xml:"result,attr"`
	Traits  []Trait  `xml:"traits>trait"`
}

type Trait struct {
	XMLName xml.Name `xml:"trait"`
	Name    string   `xml:"name,attr"`
	Value   string   `xml:"value,attr"`
}

type JUnitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type JUnitProperties struct {
	Properties []JUnitProperty `xml:"property"`
}

type JUnitTestCase struct {
	Classname  string           `xml:"classname,attr"`
	Name       string           `xml:"name,attr"`
	Time       string           `xml:"time,attr"`
	Failure    *JUnitFailure    `xml:"failure"`
	Error      *JUnitError      `xml:"error"`
	Properties *JUnitProperties `xml:"properties"`
}

type JUnitFailure struct {
	Message string `xml:"message,attr"`
	Text    string `xml:",chardata"`
}

type JUnitError struct {
	Message string `xml:"message,attr"`
	Text    string `xml:",chardata"`
}

type JUnitTestSuite struct {
	TestCases []JUnitTestCase `xml:"testcase"`
}

type JUnitTestSuites struct {
	TestSuites []JUnitTestSuite `xml:"testsuite"`
}

type XUnitProperty struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}

type XUnitProperties struct {
	Properties []XUnitProperty `xml:"property"`
}

type XUnitTest struct {
	Name       string           `xml:"name,attr"`
	Type       string           `xml:"type,attr"`
	Method     string           `xml:"method,attr"`
	Time       string           `xml:"time,attr"`
	Result     string           `xml:"result,attr"`
	Properties *XUnitProperties `xml:"properties,omitempty"`
}

type XUnitCollection struct {
	Name  string      `xml:"name,attr"`
	Time  string      `xml:"time,attr"`
	Tests []XUnitTest `xml:"test"`
}

type XUnitAssembly struct {
	Name        string            `xml:"name,attr"`
	RunDate     string            `xml:"run-date,attr"`
	RunTime     string            `xml:"run-time,attr"`
	Collections []XUnitCollection `xml:"collection"`
}

type XUnitAssemblies struct {
	XMLName    xml.Name        `xml:"assemblies"`
	Assemblies []XUnitAssembly `xml:"assembly"`
}
