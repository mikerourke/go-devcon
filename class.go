package devcon

import (
	"errors"
	"regexp"
	"strings"
)

// ClassFilterType is used to indicate the filter option for setup classes.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/develop/device-filter-driver-ordering for more information.
type ClassFilterType string

const (
	// ClassFilterUpper represent the upper-level filter drivers.
	ClassFilterUpper ClassFilterType = "upper"

	// ClassFilterLower represent the lower-level filter drivers.
	ClassFilterLower ClassFilterType = "lower"
)

var reListClassHasDevices = regexp.MustCompile(`device\(s\) for setup class "(.*)"`)

// DeviceSetupClass contains the name and description of a device setup class, which
// represent devices that are set up and configured in the same manner.
//
// For example, SCSI media changer devices are grouped into the MediumChanger
// device setup class. The device setup class defines the class installer and
// class co-installers that are involved in installing the device.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/install/overview-of-device-setup-classes for more information.
type DeviceSetupClass struct {
	// Name is the name of the class.
	Name string `json:"name"`

	// Description is a brief description of the class.
	Description string `json:"description"`
}

// ClassFilterChangeResult is returned from the ClassFilter method and contains
// details regarding the result of the operation.
type ClassFilterChangeResult struct {
	// Filters is a slice of the filters that were either changed or queried.
	Filters []string

	// RequiresReboot indicates that a reboot is required as a result of the
	// change.
	RequiresReboot bool

	// WasChanged indicates if the filter were changed or just queried.
	WasChanged bool
}

// Classes returns all device setup classes, including classes that devices on the
// system do not use.
//
// The results are returned in the order that they appear in the registry
// (alphanumeric order by GUID).
//
// Running with the WithRemoteComputer() option is allowed.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-classes for more information.
func (dc *DevCon) Classes() ([]DeviceSetupClass, error) {
	lines, err := dc.run(commandClasses)
	if err != nil {
		return nil, err
	}

	classes := make([]DeviceSetupClass, 0)

	for _, valuePair := range parseColonSeparatedLines(lines) {
		class := DeviceSetupClass{
			Name:        valuePair[0],
			Description: valuePair[1],
		}

		classes = append(classes, class)
	}

	return classes, nil
}

// ClassFilter adds, deletes, displays, and changes the order of filter drivers
// for a device setup class. Omitting the drivers parameter performs a query
// operation that doesn't make any changes.
//
// The class specifies the device setup class. The filter can be ClassFilterUpper
// to indicate that the specified drivers are upper-class filter drivers, or
// ClassFilterLower to indicate that the specified drivers are lower-class filter
// drivers.
//
// Cannot be run with the WithRemoteComputer() option.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-classfilter for more information.
func (dc *DevCon) ClassFilter(
	class string,
	filter ClassFilterType,
	drivers ...string,
) (*ClassFilterChangeResult, error) {
	args := []string{class, string(filter)}

	if len(drivers) != 0 {
		args = append(args, drivers...)
	}

	lines, err := dc.run(commandClassFilter, args...)
	if err != nil {
		return nil, err
	}

	changeResult := &ClassFilterChangeResult{
		Filters:        nil,
		RequiresReboot: false,
		WasChanged:     false,
	}

	for _, line := range lines {
		if !strings.HasPrefix(line, " ") {
			changeResult.WasChanged = !strings.Contains(line, "unchanged")
			changeResult.RequiresReboot = strings.Contains(line, "must be restarted")
		} else {
			changeResult.Filters = append(changeResult.Filters, trimSpaces(line))
		}
	}

	return changeResult, nil
}

// ListClass returns all devices in the specified device setup classes in a
// map with the key equal to the class name and value equal to a slice of the
// corresponding devices.
//
// Each entry in a setup class display represents one device. The entry consists
// of the unique instance name and a description of the device.
//
// Running with the WithRemoteComputer() option is allowed.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-listclass for more information.
func (dc *DevCon) ListClass(classes ...string) (map[string][]Device, error) {
	if len(classes) == 0 {
		return nil, errors.New("at least one class is required")
	}

	lines, err := dc.run(commandListClass, classes...)
	if err != nil {
		return nil, err
	}

	devicesByClassMap := make(map[string][]Device)

	if len(lines) == 0 {
		return devicesByClassMap, nil
	}

	groupIndices := make([]int, 0)

	for index, line := range lines {
		if strings.HasPrefix(line, "Listing") || strings.HasPrefix(line, "No devices") {
			groupIndices = append(groupIndices, index)
		}
	}

	groupIndices = append(groupIndices, len(lines))

	for index, groupStart := range groupIndices {
		nextIndex := index + 1
		if len(groupIndices) == nextIndex {
			break
		}

		groupEnd := groupIndices[nextIndex]

		var class string

		devices := make([]Device, 0)

		for lineIndex := groupStart; lineIndex < groupEnd; lineIndex++ {
			line := lines[lineIndex]

			if lineIndex == groupStart {
				matches := reListClassHasDevices.FindStringSubmatch(line)
				if matches != nil {
					class = strings.ToLower(matches[1])
				}
			} else {
				valuePair := parseColonSeparatedLine(line)

				if valuePair != nil {
					devices = append(devices, Device{
						ID:   valuePair[0],
						Name: valuePair[1],
					})
				}
			}
		}

		if class != "" {
			devicesByClassMap[class] = devices
		}
	}

	return devicesByClassMap, nil
}
