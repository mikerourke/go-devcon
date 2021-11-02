package devcon

import "strings"

type searchStatus int

const (
	searchNone searchStatus = iota
	searchHW
	searchCompat
)

type HwID struct {
	DeviceID      string   `json:"deviceId"`
	DeviceName    string   `json:"deviceName"`
	HardwareIDs   []string `json:"hardwareIds"`
	CompatibleIDs []string `json:"compatibleIds"`
}

func (dc *DevCon) HwIDs() ([]HwID, error) {
	lines, err := dc.runWithoutArgs(commandHwIDs)
	if err != nil {
		return nil, err
	}

	return parseHwIDs(lines), nil
}

func parseHwIDs(lines []string) []HwID {
	groupIndices := make([]int, 0)

	for index, line := range lines {
		if !strings.HasPrefix(line, " ") {
			groupIndices = append(groupIndices, index)
		}
	}

	hwids := make([]HwID, 0)

	for index, groupStart := range groupIndices {
		nextIndex := index + 1
		if len(groupIndices) == nextIndex {
			break
		}

		groupEnd := groupIndices[nextIndex]

		hwid := HwID{
			HardwareIDs:   make([]string, 0),
			CompatibleIDs: make([]string, 0),
		}
		status := searchNone

		for lineIndex := groupStart; lineIndex < groupEnd; lineIndex++ {
			thisLine := lines[lineIndex]

			if lineIndex == groupStart {
				hwid.DeviceID = thisLine

				status = searchNone
			} else if lineIndex == groupStart+1 {
				nameParams := parseParams(reName, thisLine)

				if name, ok := nameParams["Name"]; ok {
					hwid.DeviceName = name
				}
			} else if strings.Contains(thisLine, "Hardware ID") {
				status = searchHW
			} else if strings.Contains(thisLine, "Compatible ID") {
				status = searchCompat
			} else {
				idLine := strings.Trim(thisLine, " ")

				if status == searchHW {
					hwid.HardwareIDs = append(hwid.HardwareIDs, idLine)
				} else if status == searchCompat {
					hwid.CompatibleIDs = append(hwid.CompatibleIDs, idLine)
				}
			}

			if lineIndex == groupEnd-1 && hwid.DeviceName != "" {
				hwids = append(hwids, hwid)
			}
		}
	}

	return hwids
}
