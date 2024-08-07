package structs

type JUnitTestSuites struct {
	TestSuites []JUnitTestSuite `xml:"testsuite"`
}

type JUnitTestSuite struct {
	ID         string          `xml:"id,attr"`
	Name       string          `xml:"name,attr"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Errors     int             `xml:"errors,attr"`
	Skipped    int             `xml:"skipped,attr"`
	Time       float64         `xml:"time,attr"`
	TestCases  []JUnitTestCase `xml:"testcase"`
	Properties []Property      `xml:"properties>property"`
}

type JUnitTestCase struct {
	ID         string     `xml:"id,attr"`
	ClassName  string     `xml:"classname,attr"`
	Name       string     `xml:"name,attr"`
	Time       float64    `xml:"time,attr"`
	Properties []Property `xml:"properties>property"`
}

type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr"`
}
