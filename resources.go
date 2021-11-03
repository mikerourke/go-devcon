package devcon

import (
	"strconv"
	"strings"
)

type DeviceResource struct {
	Name  string      `json:"name"`
	Value interface{} `json:"value"`
}

type DeviceResourceUsage struct {
	DeviceID   string           `json:"deviceId"`
	DeviceName string           `json:"deviceName"`
	Resources  []DeviceResource `json:"resources"`
}

func (dc *DevCon) Resources() ([]DeviceResourceUsage, error) {
	lines, err := dc.run(commandResources)
	if err != nil {
		return nil, err
	}

	return parseResources(lines), nil
}

func parseResources(lines []string) []DeviceResourceUsage {
	groupIndices := make([]int, 0)

	for index, line := range lines {
		if !strings.HasPrefix(line, " ") {
			groupIndices = append(groupIndices, index)
		}
	}

	drus := make([]DeviceResourceUsage, 0)

	for index, groupStart := range groupIndices {
		nextIndex := index + 1
		if len(groupIndices) == nextIndex {
			break
		}

		groupEnd := groupIndices[nextIndex]

		dru := DeviceResourceUsage{
			Resources: make([]DeviceResource, 0),
		}

		for lineIndex := groupStart; lineIndex < groupEnd; lineIndex++ {
			thisLine := strings.Trim(lines[lineIndex], " ")

			if lineIndex == groupStart {
				dru.DeviceID = thisLine
			} else if lineIndex == groupStart+1 {
				nameParams := parseParams(reName, thisLine)

				if name, ok := nameParams["Name"]; ok {
					dru.DeviceName = name
				}
			} else {
				params := parseParams(reResource, thisLine)

				dr := DeviceResource{}
				if name, ok := params["Name"]; ok {
					dr.Name = strings.Trim(name, " ")
				}

				if value, ok := params["Value"]; ok {
					value = strings.Trim(value, " ")
					number, err := strconv.Atoi(value)
					if err == nil {
						dr.Value = number
					} else {
						dr.Value = value
					}
				}

				if dr.Name != "" {
					dru.Resources = append(dru.Resources, dr)
				}
			}

			if lineIndex == groupEnd-1 && dru.DeviceName != "" {
				drus = append(drus, dru)
			}
		}
	}

	return drus
}
