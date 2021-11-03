package devcon

import "strings"

type DriverStatus struct {
	DeviceID   string `json:"deviceId"`
	DeviceName string `json:"deviceName"`
	IsRunning  bool   `json:"isRunning"`
}

func (dc *DevCon) Status() ([]DriverStatus, error) {
	lines, err := dc.run(commandStatus)
	if err != nil {
		return nil, err
	}

	return parseStatus(lines), nil
}

func parseStatus(lines []string) []DriverStatus {
	groupIndices := make([]int, 0)

	for index, line := range lines {
		if !strings.HasPrefix(line, " ") {
			groupIndices = append(groupIndices, index)
		}
	}

	statuses := make([]DriverStatus, 0)

	for index, groupStart := range groupIndices {
		nextIndex := index + 1
		if len(groupIndices) == nextIndex {
			break
		}

		groupEnd := groupIndices[nextIndex]

		status := DriverStatus{}

		for lineIndex := groupStart; lineIndex < groupEnd; lineIndex++ {
			thisLine := lines[lineIndex]

			if lineIndex == groupStart {
				status.DeviceID = thisLine
			} else if lineIndex == groupStart+1 {
				nameParams := parseParams(reName, thisLine)

				if name, ok := nameParams["Name"]; ok {
					status.DeviceName = name
				}
			} else {
				statusLine := strings.Trim(thisLine, " ")
				status.IsRunning = statusLine == "Driver is running."
			}

			if lineIndex == groupEnd-1 && status.DeviceName != "" {
				statuses = append(statuses, status)
			}
		}
	}

	return statuses
}
