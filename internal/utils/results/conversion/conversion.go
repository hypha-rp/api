package conversion

import (
	"encoding/xml"
	"fmt"
	"hypha/api/internal/utils/results/structs"
	"io/ioutil"
)

func ConvertJUnitToXUnit(junitFile string) (string, error) {
	junitData, err := ioutil.ReadFile(junitFile)
	if err != nil {
		return "", fmt.Errorf("error reading JUnit file: %v", err)
	}

	var junitTestSuites structs.JUnitTestSuites
	err = xml.Unmarshal(junitData, &junitTestSuites)
	if err != nil {
		return "", fmt.Errorf("error unmarshalling JUnit XML: %v", err)
	}

	xunitAssemblies := structs.XUnitAssemblies{
		Assemblies: []structs.XUnitAssembly{
			{
				Name:    "JUnitConverted",
				RunDate: "2023-10-01",
				RunTime: "00:00:00",
			},
		},
	}

	for _, suite := range junitTestSuites.TestSuites {
		for _, testcase := range suite.TestCases {
			result := "passed"
			if testcase.Failure != nil {
				result = "failed"
			} else if testcase.Error != nil {
				result = "error"
			}

			var xunitProperties *structs.XUnitProperties
			if testcase.Properties != nil {
				xunitProperties = &structs.XUnitProperties{}
				for _, prop := range testcase.Properties.Properties {
					xunitProperties.Properties = append(xunitProperties.Properties, structs.XUnitProperty{
						Name:  prop.Name,
						Value: prop.Value,
					})
				}
			}

			xunitTest := structs.XUnitTest{
				Name:       testcase.Name,
				Type:       testcase.Classname,
				Method:     testcase.Name,
				Time:       testcase.Time,
				Result:     result,
				Properties: xunitProperties,
			}

			xunitCollection := structs.XUnitCollection{
				Name:  testcase.Classname,
				Time:  testcase.Time,
				Tests: []structs.XUnitTest{xunitTest},
			}

			xunitAssemblies.Assemblies[0].Collections = append(xunitAssemblies.Assemblies[0].Collections, xunitCollection)
		}
	}

	xunitData, err := xml.MarshalIndent(xunitAssemblies, "", "  ")
	if err != nil {
		return "", fmt.Errorf("error marshalling XUnit XML: %v", err)
	}

	xunitFile := "xunit-results.xml"
	err = ioutil.WriteFile(xunitFile, xunitData, 0644)
	if err != nil {
		return "", fmt.Errorf("error writing XUnit file: %v", err)
	}

	return xunitFile, nil
}
