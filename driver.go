package devcon

import (
	"fmt"
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

// DriverFileGroup describes the INF file and section, as well as associated
// files for a device.
type DriverFileGroup struct {
	// Device is the device with which the files are associated.
	Device Device `json:"device"`

	// INFFile is the driver file which is used by the device.
	INFFile string `json:"infFile"`

	// INFSection is the corresponding section of the INF file.
	// TODO: This is a bad description.
	INFSection string `json:"infSection"`

	// Files are the driver files associated with the device driver.
	Files []string `json:"files"`
}

// DriverNode describes the components of a driver package.
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/install/components-of-a-driver-package
// for more information about driver packages.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/install/how-setup-ranks-drivers--windows-vista-and-later-
// for more information about node rank.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/debugger/device-node-status-flags
// for more information about node flags.
type DriverNode struct {
	// NodeNumber represents the node order for the device driver.
	NodeNumber int `json:"nodeNumber"`

	// INFFile is the fully qualified path to the INF file.
	INFFile string `json:"infFile"`

	// INFSection is the section of the INF file to which this device
	// corresponds.
	INFSection string `json:"infSection"`

	// Description is the description of the device.
	Description string `json:"description"`

	// Manufacturer is the name of the device manufacturer.
	Manufacturer string `json:"manufacturer"`

	// Provider is the name of the driver provider (e.g. Microsoft).
	Provider string `json:"provider"`

	// Date is the date associated with the current driver version.
	Date string `json:"date"`

	// Version is the current version of the driver.
	Version string `json:"version"`

	// NodeRank indicates how well the driver matches the device. A driver rank
	// is represented by an integer that is equal to or greater than zero. The
	// lower the rank, the better a match the driver is for the device.
	NodeRank int `json:"nodeRank"`

	// NodeFlags describe the status of a device.
	NodeFlags int `json:"nodeFlags"`

	// IsDigitallySigned indicates that the driver has been digitally signed.
	IsDigitallySigned bool `json:"isDigitallySigned"`
}

// DriverNodeGroup contains device details as well as details for the corresponding
// driver nodes.
type DriverNodeGroup struct {
	// Device is the device with which the nodes are associated.
	Device Device `json:"device"`

	// Nodes are DriverNode records that describe the nodes associated with the
	// device driver.
	Nodes []DriverNode `json:"nodes"`
}

// DriverFiles returns the full path and file name of installed INF files and
// device driver files for the specified devices.
//
// Cannot be run with the WithRemoteComputer() option.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-driverfiles for more information.
func (dc *DevCon) DriverFiles(idOrClass string, idsOrClasses ...string) ([]DriverFileGroup, error) {
	allIdsOrClasses := concatIdsOrClasses(idOrClass, idsOrClasses...)

	lines, err := dc.run(commandDriverFiles, allIdsOrClasses...)
	if err != nil {
		return nil, err
	}

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
			Device: Device{
				ID:   "",
				Name: "",
			},
			INFFile:    "",
			INFSection: "",
			Files:      make([]string, 0),
		}

		for lineIndex := groupStart; lineIndex < groupEnd; lineIndex++ {
			line := lines[lineIndex]

			switch {
			case lineIndex == groupStart:
				fileGroup.Device.ID = line

			case lineIndex == groupStart+1:
				nameParams := parseParams(reName, line)

				if name, ok := nameParams["Name"]; ok {
					fileGroup.Device.Name = name
				}

			case reDriverNoInfo.MatchString(line):
				continue

			default:
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

	return fileGroups, nil
}

// DriverNodes returns all driver packages that are compatible with the device,
// along with their version and ranking.
//
// The DriverNodes method is particularly useful for troubleshooting setup
// problems. For example, you can use it to determine whether a Windows INF
// file or a customized third-party INF file was used for a device.
//
// Cannot be run with the WithRemoteComputer() option.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-drivernodes for more information.
//
//nolint:funlen // Function is long, but simple.
func (dc *DevCon) DriverNodes(idOrClass string, idsOrClasses ...string) ([]DriverNodeGroup, error) {
	allIdsOrClasses := concatIdsOrClasses(idOrClass, idsOrClasses...)

	lines, err := dc.run(commandDriverNodes, allIdsOrClasses...)
	if err != nil {
		return nil, err
	}

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
			Device: Device{
				ID:   "",
				Name: "",
			},
			Nodes: make([]DriverNode, 0),
		}
		node := DriverNode{
			NodeNumber:        0,
			INFFile:           "",
			INFSection:        "",
			Description:       "",
			Manufacturer:      "",
			Provider:          "",
			Date:              "",
			Version:           "",
			NodeRank:          0,
			NodeFlags:         0,
			IsDigitallySigned: false,
		}

		for lineIndex := groupStart; lineIndex < groupEnd; lineIndex++ {
			line := lines[lineIndex]

			switch {
			case lineIndex == groupStart:
				nodeGroup.Device.ID = line

			case lineIndex == groupStart+1:
				nameParams := parseParams(reName, line)

				if name, ok := nameParams["Name"]; ok {
					nodeGroup.Device.Name = name
				}

			case reHash.MatchString(line):
				matches := reDriverNode.FindStringSubmatch(line)

				if matches != nil {
					number, _ := strconv.Atoi(matches[1])
					node.NodeNumber = number
				}

			case strings.Contains(line, "No DriverNodes"):
				continue

			default:
				params := parseParams(reFieldIsValue, line)

				if field, ok := params["Field"]; ok {
					value, ok := params["Value"]
					if !ok {
						continue
					}

					switch trimSpaces(field) {
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

						node = DriverNode{
							NodeNumber:        0,
							INFFile:           "",
							INFSection:        "",
							Description:       "",
							Manufacturer:      "",
							Provider:          "",
							Date:              "",
							Version:           "",
							NodeRank:          0,
							NodeFlags:         0,
							IsDigitallySigned: false,
						}
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

	return nodeGroups, nil
}

// Update forcibly replaces the current device drivers for a specified device
// with drivers listed in the specified INF file.
//
// Update forces an update to the most appropriate drivers in the specified INF
// file, even if those drivers are older or less appropriate than the current drivers
// or the drivers in a different INF file.
//
// You cannot use Update to update drivers for non-present devices.
//
// Before updating the driver for any device, determine which devices will
// be affected. To do so, pass the name to the HwIDs() function:
//
//	dc.HwIDs("ISAPNP\CSC4324\0")
//
// Or with the DriverFiles() function:
//
//	dc.DriverFiles("ISAPNP\CSC4324\0")
//
// The system might need to be rebooted to make this change effective. To reboot
// the system if required, use:
//
//	dc.WithConditionalReboot().Update()
//
// Cannot be run with the WithRemoteComputer() option.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-update for more information.
func (dc *DevCon) Update(infFile string, hardwareID string) error {
	lines, err := dc.run(commandUpdate, infFile, hardwareID)

	if substrInLines(lines, "success") == -1 {
		return fmt.Errorf("error updating driver: %s", strings.Join(lines, ". "))
	}

	return err
}

// UpdateNI forcibly replaces the current device drivers with drivers listed in
// the specified INF file without prompting the user for information or
// confirmation.
//
// This method suppresses all user prompts that require a response and assumes
// the default response. As a result, you cannot use this operation to install
// unsigned drivers. To display user prompts during an update, use Update().
//
// This method forces an update, even if the drivers in the specified INF file
// are older or less appropriate than the current drivers.
//
// Before updating the driver for any device, determine which devices will
// be affected. To do so, pass the name to the HwIDs() function:
//
//	dc.HwIDs("ISAPNP\CSC4324\0")
//
// Or with the DriverFiles() function:
//
//	dc.DriverFiles("ISAPNP\CSC4324\0")
//
// The system might need to be rebooted to make this change effective. To reboot
// the system if required, use:
//
//	dc.WithConditionalReboot().UpdateNI()
//
// Cannot be run with the WithRemoteComputer() option.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-updateni for more information.
func (dc *DevCon) UpdateNI(infFile string, hardwareID string) error {
	lines, err := dc.run(commandUpdateNI, infFile, hardwareID)

	if substrInLines(lines, "success") == -1 {
		return fmt.Errorf("error updating driver: %s", strings.Join(lines, ". "))
	}

	return err
}
