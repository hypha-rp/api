package results

type JUnitTestSuites struct {
	Name       string           `xml:"name,attr"`
	Tests      int              `xml:"tests,attr"`
	Failures   int              `xml:"failures,attr"`
	Errors     int              `xml:"errors,attr"`
	Skipped    int              `xml:"skipped,attr"`
	Assertions int              `xml:"assertions,attr"`
	Time       float64          `xml:"time,attr"`
	Timestamp  string           `xml:"timestamp,attr"`
	TestSuites []JUnitTestSuite `xml:"testsuite"`
}

type JUnitTestSuite struct {
	ID         string          `xml:"id,attr"`
	Name       string          `xml:"name,attr"`
	Tests      int             `xml:"tests,attr"`
	Failures   int             `xml:"failures,attr"`
	Errors     int             `xml:"errors,attr"`
	Skipped    int             `xml:"skipped,attr"`
	Assertions int             `xml:"assertions,attr"`
	Time       float64         `xml:"time,attr"`
	File       string          `xml:"file,attr"`
	TestCases  []JUnitTestCase `xml:"testcase"`
	Properties []Property      `xml:"properties>property"`
	SystemOut  string          `xml:"system-out,omitempty"`
	SystemErr  string          `xml:"system-err,omitempty"`
}

type JUnitTestCase struct {
	ID         string     `xml:"id,attr"`
	ClassName  string     `xml:"classname,attr"`
	Name       string     `xml:"name,attr"`
	Time       float64    `xml:"time,attr"`
	Assertions int        `xml:"assertions,attr"`
	File       string     `xml:"file,attr"`
	Line       int        `xml:"line,attr"`
	Status     string     `xml:"-"`
	Failure    *Failure   `xml:"failure,omitempty"`
	Error      *Error     `xml:"error,omitempty"`
	Skipped    *Skipped   `xml:"skipped,omitempty"`
	Properties []Property `xml:"properties>property"`
	SystemOut  string     `xml:"system-out,omitempty"`
	SystemErr  string     `xml:"system-err,omitempty"`
}

type Failure struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
}

type Error struct {
	Message string `xml:"message,attr"`
	Type    string `xml:"type,attr"`
}

type Skipped struct {
	Message string `xml:"message,attr"`
}

type Property struct {
	Name  string `xml:"name,attr"`
	Value string `xml:"value,attr,omitempty"`
	Text  string `xml:",chardata"`
}