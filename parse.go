package devcon

import (
	"regexp"
	"strings"
)

var (
	reName     = regexp.MustCompile(`Name: (?P<Name>.*)`)
	reKeyValue = regexp.MustCompile(`(?P<Key>.*): (?P<Value>.*)`)
)

// parseColonSeparatedLines returns a slice of 2-element slices from the
// parseColonSeparatedLine function where each slice corresponds to a line.
func parseColonSeparatedLines(lines []string) [][]string {
	values := make([][]string, 0)

	const ValidPairCount = 2

	for _, line := range lines {
		valuePairs := parseColonSeparatedLine(line)
		if len(valuePairs) == ValidPairCount {
			values = append(values, valuePairs)
		}
	}

	return values
}

// parseColonSeparatedLine returns a 2-element slice where the first element is
// the key (before the colon) and second is the value (after the colon).
func parseColonSeparatedLine(line string) []string {
	if len(line) == 0 {
		return nil
	}

	params := parseParams(reKeyValue, line)
	key, ok := params["Key"]

	if !ok {
		return nil
	}

	key = trimSpaces(key)

	value, ok := params["Value"]
	if ok {
		value = trimSpaces(value)
	}

	return []string{key, value}
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
	lines := strings.Split(contents, "\n")

	validLines := make([]string, 0)

	for _, line := range lines {
		if strings.Contains(line, "\r") {
			//nolint:gocritic // ReplaceAll not supported in Go 1.10.7.
			validLines = append(validLines, strings.Replace(line, "\r", "", -1))
		} else {
			validLines = append(validLines, line)
		}
	}

	return validLines
}

// trimSpaces removes any surrounding whitespace from the specified value.
func trimSpaces(value string) string {
	return strings.Trim(value, " ")
}
