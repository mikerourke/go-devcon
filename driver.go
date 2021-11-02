package devcon

import (
	"strconv"
	"strings"
)

type DriverFileGroup struct {
	DeviceID   string   `json:"deviceId"`
	DeviceName string   `json:"deviceName"`
	INFFile    string   `json:"infFile"`
	INFSection string   `json:"infSection"`
	Files      []string `json:"files"`
}

type DriverNode struct {
	NodeNumber        int    `json:"nodeNumber"`
	INFFile           string `json:"infFile"`
	INFSection        string `json:"infSection"`
	Description       string `json:"description"`
	Manufacturer      string `json:"manufacturer"`
	Provider          string `json:"provider"`
	Date              string `json:"date"`
	Version           string `json:"version"`
	NodeRank          int    `json:"nodeRank"`
	NodeFlags         int    `json:"nodeFlags"`
	IsDigitallySigned bool   `json:"isDigitallySigned"`
}

type DriverNodeGroup struct {
	DeviceID   string       `json:"deviceId"`
	DeviceName string       `json:"deviceName"`
	Nodes      []DriverNode `json:"nodes"`
}

func (dc *DevCon) DriverFiles() ([]DriverFileGroup, error) {
	lines, err := dc.runWithoutArgs(commandDriverFiles)
	if err != nil {
		return nil, err
	}

	return parseDriverFileGroups(lines), nil
}

func (dc *DevCon) DriverNodes() ([]DriverNodeGroup, error) {
	lines, err := dc.runWithoutArgs(commandDriverNodes)
	if err != nil {
		return nil, err
	}

	return parseDriverNodeGroups(lines), nil
}

func parseDriverFileGroups(lines []string) []DriverFileGroup {
	driverGroupLines := make([][]string, 0)

	driverLines := make([]string, 0)

	for index, line := range lines {
		if !strings.HasPrefix(line, " ") {
			if index != 0 {
				driverGroupLines = append(driverGroupLines, driverLines)

				driverLines = make([]string, 0)
			}
		}

		driverLines = append(driverLines, line)
	}

	fileGroups := make([]DriverFileGroup, 0)

	for _, lines = range driverGroupLines {
		fileGroup := parseDriverFileGroup(lines)

		if fileGroup.DeviceName != "" {
			fileGroups = append(fileGroups, fileGroup)
		}
	}

	return fileGroups
}

func parseDriverFileGroup(lines []string) DriverFileGroup {
	fileGroup := DriverFileGroup{
		DeviceID:   "",
		DeviceName: "",
		INFFile:    "",
		INFSection: "",
		Files:      make([]string, 0),
	}

	for _, line := range lines {
		if !strings.HasPrefix(line, " ") {
			fileGroup.DeviceID = line
		}

		nameParams := parseParams(reName, line)
		if name, ok := nameParams["Name"]; ok && name != "" {
			fileGroup.DeviceName = name
		}

		if reDriverNoInfo.MatchString(line) {
			continue
		}

		driverParams := parseParams(reDriverInstalled, line)
		if infPath, ok := driverParams["INFFile"]; ok && infPath != "" {
			fileGroup.INFFile = infPath
		}

		if infName, ok := driverParams["INFSection"]; ok && infName != "" {
			fileGroup.INFSection = infName
		}

		fileResult := reDriverFilePath.FindStringSubmatch(line)
		if fileResult != nil {
			fileGroup.Files = append(fileGroup.Files, fileResult[0])
		}
	}

	return fileGroup
}

func parseDriverNodeGroups(lines []string) []DriverNodeGroup {
	groupIndices := make([]int, 0)

	for index, line := range lines {
		if !strings.HasPrefix(line, " ") && !reHash.MatchString(line) {
			groupIndices = append(groupIndices, index)
		}
	}

	nodeGroups := make([]DriverNodeGroup, 0)

	for index, groupStart := range groupIndices {
		nextIndex := index + 1
		if len(groupIndices) == nextIndex {
			break
		}

		groupEnd := groupIndices[nextIndex]

		nodeGroup := DriverNodeGroup{
			Nodes: make([]DriverNode, 0),
		}

		node := DriverNode{}

		for lineIndex := groupStart; lineIndex < groupEnd; lineIndex++ {
			thisLine := lines[lineIndex]

			if lineIndex == groupStart {
				nodeGroup.DeviceID = thisLine
			} else if lineIndex == groupStart+1 {
				nameParams := parseParams(reName, thisLine)

				if name, ok := nameParams["Name"]; ok {
					nodeGroup.DeviceName = name
				}
			} else if reHash.MatchString(thisLine) {
				matches := reDriverNode.FindStringSubmatch(thisLine)

				if matches != nil {
					number, _ := strconv.Atoi(matches[1])
					node.NodeNumber = number
				}
			} else if strings.Contains(thisLine, "No DriverNodes") {
				continue
			} else {
				params := parseParams(reFieldIsValue, thisLine)
				if field, ok := params["Field"]; ok {
					value, ok := params["Value"]
					if !ok {
						continue
					}

					switch strings.Trim(field, " ") {
					case "Inf file":
						node.INFFile = value

					case "Inf section":
						node.INFSection = value

					case "Driver description":
						node.Description = value

					case "Manufacturer name":
						node.Manufacturer = value

					case "Provider name":
						node.Provider = value

					case "Driver date":
						node.Date = value

					case "Driver version":
						node.Version = value

					case "Driver node rank":
						number, _ := strconv.Atoi(value)
						node.NodeRank = number
					}

					if value == "digitally signed" {
						node.IsDigitallySigned = true
						nodeGroup.Nodes = append(nodeGroup.Nodes, node)

						node = DriverNode{}
					}
				}

				params = parseParams(reFieldAreValue, thisLine)
				if value, ok := params["Value"]; ok {
					number, _ := strconv.Atoi(value)
					node.NodeFlags = number
				}
			}

			if lineIndex == groupEnd-1 && nodeGroup.DeviceName != "" {
				nodeGroups = append(nodeGroups, nodeGroup)
			}
		}
	}

	return nodeGroups
}
