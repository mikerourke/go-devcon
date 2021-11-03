package devcon

import "strings"

type stackSearchStatus int

const (
	stackSearchNone stackSearchStatus = iota
	stackSearchUpperFilter
	stackSearchService
	stackSearchLowerFilter
)

type DriverStack struct {
	DeviceID           string `json:"deviceId"`
	DeviceName         string `json:"deviceName"`
	SetupClassGUID     string `json:"setupClassGUID"`
	SetupClassName     string `json:"setupClassName"`
	ControllingService string `json:"controllingService"`
	UpperFilters       string `json:"upperFilters"`
	LowerFilters       string `json:"lowerFilters"`
}

func (dc *DevCon) Stack() ([]DriverStack, error) {
	lines, err := dc.run(commandStack)
	if err != nil {
		return nil, err
	}

	return parseStacks(lines), nil
}

func parseStacks(lines []string) []DriverStack {
	groupIndices := make([]int, 0)

	for index, line := range lines {
		if !strings.HasPrefix(line, " ") {
			groupIndices = append(groupIndices, index)
		}
	}

	stacks := make([]DriverStack, 0)

	for index, groupStart := range groupIndices {
		nextIndex := index + 1
		if len(groupIndices) == nextIndex {
			break
		}

		groupEnd := groupIndices[nextIndex]

		stack := DriverStack{}

		search := stackSearchNone

		for lineIndex := groupStart; lineIndex < groupEnd; lineIndex++ {
			thisLine := strings.Trim(lines[lineIndex], " ")

			if lineIndex == groupStart {
				stack.DeviceID = thisLine
			} else if lineIndex == groupStart+1 {
				nameParams := parseParams(reName, thisLine)

				if name, ok := nameParams["Name"]; ok {
					stack.DeviceName = name
				}
			} else if strings.Contains(thisLine, "Setup Class") {
				params := parseParams(reSetupClass, thisLine)

				if guid, ok := params["GUID"]; ok {
					stack.SetupClassGUID = guid
				}

				if name, ok := params["Name"]; ok {
					stack.SetupClassName = name
				}
			} else if strings.Contains(thisLine, "pper filters:") {
				search = stackSearchUpperFilter
			} else if strings.Contains(thisLine, "service:") {
				search = stackSearchService
			} else if strings.Contains(thisLine, "ower filters:") {
				search = stackSearchLowerFilter
			} else {
				if search == stackSearchUpperFilter {
					stack.UpperFilters = thisLine
				} else if search == stackSearchService {
					stack.ControllingService = thisLine
				} else if search == stackSearchLowerFilter {
					stack.LowerFilters = thisLine
				}
			}

			if lineIndex == groupEnd-1 && stack.DeviceName != "" {
				stacks = append(stacks, stack)
			}
		}
	}

	return stacks
}
