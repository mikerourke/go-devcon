package devcon

import (
	"regexp"
	"strconv"
	"strings"
)

var (
	reDriverFilePath  = regexp.MustCompile(`C:\\(.*)`)
	reDriverInstalled = regexp.MustCompile(`Driver installed from (?P<INFFile>.*) \[(?P<INFSection>.*)].*`)
	reDriverNode      = regexp.MustCompile(`DriverNode #(.*):`)
	reDriverNoInfo    = regexp.MustCompile(`No driver information`)
	reFieldAreValue   = regexp.MustCompile(`(?P<Field>.*) are (?P<Value>.*)`)
	reFieldIsValue    = regexp.MustCompile(`(?P<Field>.*) is (?P<Value>.*)`)
	reHash            = regexp.MustCompile(`#`)
)

type DriverFileGroup struct {
	Device     Device   `json:"device"`
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
	Device Device       `json:"device"`
	Nodes  []DriverNode `json:"nodes"`
}

// DriverFiles returns the full path and file name of installed INF files and
// device driver files for the specified devices. Valid only on the local computer.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-driverfiles for more information.
func (dc *DevCon) DriverFiles() ([]DriverFileGroup, error) {
	lines, err := dc.run(commandDriverFiles)
	if err != nil {
		return nil, err
	}

	return parseDriverFileGroups(lines), nil
}

func (dc *DevCon) DriverNodes() ([]DriverNodeGroup, error) {
	lines, err := dc.run(commandDriverNodes)
	if err != nil {
		return nil, err
	}

	return parseDriverNodeGroups(lines), nil
}

// Update forcibly replaces the current device drivers for a specified device
// with drivers listed in the specified INF file. Valid only on the local computer.
//
// Notes
// Update forces an update to the most appropriate drivers in the specified INF
// file, even if those drivers are older or less appropriate than the current drivers
// or the drivers in a different INF file.
//
// You cannot use Update to update drivers for non-present devices.
//
// Before updating the driver for any device, determine which devices will
// be affected. To do so, pass the name to the HwIDs() function:
//	dc.HwIDs("ISAPNP\CSC4324\0")
// Or with the DriverFiles() function:
//	dc.DriverFiles("ISAPNP\CSC4324\0")
//
// The system might need to be rebooted to make this change effective. To
// reboot the system, add ConditionalReboot() before Update().
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-update for more information.
func (dc *DevCon) Update(infFile string, hardwareID string) error {
	lines, err := dc.run(commandUpdate, infFile, hardwareID)

	// TODO: Parse
	dc.printResults(lines)

	return err
}

// UpdateNI forcibly replaces the current device drivers with drivers listed in
// the specified INF file without prompting the user for information or
// confirmation. Valid only on the local computer.
//
// Notes
// UpdateNI suppresses all user prompts that require a response and assumes the
// default response. As a result, you cannot use this operation to install
// unsigned drivers. To display user prompts during an update, use Update().
//
// UpdateNI forces an update, even if the drivers in the specified INF file are
// older or less appropriate than the current drivers.
//
// Before updating the driver for any device, determine which devices will
// be affected. To do so, pass the name to the HwIDs() function:
//	dc.HwIDs("ISAPNP\CSC4324\0")
// Or with the DriverFiles() function:
//	dc.DriverFiles("ISAPNP\CSC4324\0")
//
// The system might need to be rebooted to make this change effective. To
// reboot the system, add ConditionalReboot() before UpdateNI().
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-updateni for more information.
func (dc *DevCon) UpdateNI(infFile string, hardwareID string) error {
	lines, err := dc.run(commandUpdateNI, infFile, hardwareID)

	// TODO: Parse
	dc.printResults(lines)

	return err
}

func parseDriverFileGroups(lines []string) []DriverFileGroup {
	groupIndices := make([]int, 0)

	for index, line := range lines {
		if !strings.HasPrefix(line, " ") {
			groupIndices = append(groupIndices, index)
		}
	}

	groupIndices = append(groupIndices, len(lines))

	fileGroups := make([]DriverFileGroup, 0)

	for index, groupStart := range groupIndices {
		nextIndex := index + 1
		if len(groupIndices) == nextIndex {
			break
		}

		groupEnd := groupIndices[nextIndex]

		fileGroup := DriverFileGroup{
			Files: make([]string, 0),
		}

		for lineIndex := groupStart; lineIndex < groupEnd; lineIndex++ {
			line := lines[lineIndex]

			if lineIndex == groupStart {
				fileGroup.Device.ID = line
			} else if lineIndex == groupStart+1 {
				nameParams := parseParams(reName, line)

				if name, ok := nameParams["Name"]; ok {
					fileGroup.Device.Name = name
				}
			} else {
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
		}

		if fileGroup.Device.Name != "" {
			fileGroups = append(fileGroups, fileGroup)
		}
	}

	return fileGroups
}

func parseDriverNodeGroups(lines []string) []DriverNodeGroup {
	groupIndices := make([]int, 0)

	for index, line := range lines {
		if !strings.HasPrefix(line, " ") && !reHash.MatchString(line) {
			groupIndices = append(groupIndices, index)
		}
	}

	groupIndices = append(groupIndices, len(lines))

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
			line := lines[lineIndex]

			if lineIndex == groupStart {
				nodeGroup.Device.ID = line
			} else if lineIndex == groupStart+1 {
				nameParams := parseParams(reName, line)

				if name, ok := nameParams["Name"]; ok {
					nodeGroup.Device.Name = name
				}
			} else if reHash.MatchString(line) {
				matches := reDriverNode.FindStringSubmatch(line)

				if matches != nil {
					number, _ := strconv.Atoi(matches[1])
					node.NodeNumber = number
				}
			} else if strings.Contains(line, "No DriverNodes") {
				continue
			} else {
				params := parseParams(reFieldIsValue, line)
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

				params = parseParams(reFieldAreValue, line)
				if value, ok := params["Value"]; ok {
					number, _ := strconv.Atoi(value)
					node.NodeFlags = number
				}
			}

			if lineIndex == groupEnd-1 && nodeGroup.Device.Name != "" {
				nodeGroups = append(nodeGroups, nodeGroup)
			}
		}
	}

	return nodeGroups
}
