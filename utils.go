package devcon

import (
	"bytes"
	"encoding/json"
	"regexp"
	"strings"
)

// MarshalJSON returns the JSON representation of the specified interface
// without HTML escaped.
func MarshalJSON(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}

	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(t)

	return buffer.Bytes(), err
}

func parseColonSeparatedLines(lines []string) [][]string {
	values := make([][]string, 0)

	for _, line := range lines {
		if len(line) == 0 {
			continue
		}

		params := parseParams(reNameDesc, line)
		name, ok := params["Name"]

		if !ok {
			continue
		}

		name = strings.Trim(name, " ")

		desc, ok := params["Desc"]
		if ok {
			desc = strings.Trim(desc, " ")
		}

		values = append(values, []string{name, desc})
	}

	return values
}

// parseParams applies the specified regEx to the specified contents and returns
// a map of matches in the named capture groups.
func parseParams(regEx *regexp.Regexp, contents string) map[string]string {
	match := regEx.FindStringSubmatch(contents)

	paramsMap := make(map[string]string)
	for i, name := range regEx.SubexpNames() {
		if i > 0 && i <= len(match) {
			paramsMap[name] = match[i]
		}
	}

	return paramsMap
}

// splitLines splits the specified contents into lines separated by line breaks.
func splitLines(contents string) []string {
	return strings.Split(contents, "\r\n")
}
