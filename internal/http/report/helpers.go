package report

import "bytes"

// containsTestsuitesTag checks if the given XML content contains a <testsuites> tag.
// Parameters:
// - xmlContent: The XML content to check as a byte slice.
// Returns:
// - A boolean indicating whether the <testsuites> tag is present in the XML content.
func containsTestsuitesTag(xmlContent []byte) bool {
	return bytes.Contains(xmlContent, []byte("<testsuites"))
}

// wrapInTestsuitesTag wraps the given XML content in a <testsuites> tag.
// Parameters:
// - xmlContent: The XML content to wrap as a byte slice.
// Returns:
// - A new byte slice with the XML content wrapped in a <testsuites> tag.
func wrapInTestsuitesTag(xmlContent []byte) []byte {
	return append([]byte("<testsuites>"), append(xmlContent, []byte("</testsuites>")...)...)
}
