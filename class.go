package devcon

import (
	"regexp"
	"strings"
)

// ClassFilter is used to indicate the filter option for setup classes.
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/develop/device-filter-driver-ordering for more information.
type ClassFilter string

const (
	// ClassFilterUpper represent the upper-level filter drivers.
	ClassFilterUpper ClassFilter = "upper"

	// ClassFilterLower represent the lower-level filter drivers.
	ClassFilterLower ClassFilter = "lower"
)

var (
	reListClassHasDevices = regexp.MustCompile(`device\(s\) for setup class "(.*)"`)
)

// Class contains the name and description of a device setup class, which
// represent devices that are set up and configured in the same manner.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/install/overview-of-device-setup-classes for more information.
type Class struct {
	// Name is the name of the class.
	Name string `json:"name"`

	// Description is a brief description of the class.
	Description string `json:"description"`
}

// Classes returns all device setup classes, including classes that devices on the
// system do not use. Valid on local and remote computers.
//
// Notes
// Classes are returned in the order that they appear in the registry (alphanumeric
// order by GUID).
//
// To find the devices in a setup class, use the ListClass() function. To find
// the setup class of a particular device, use the Stack() function.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-classes for more information.
func (dc *DevCon) Classes() ([]Class, error) {
	lines, err := dc.run(commandClasses)
	if err != nil {
		return nil, err
	}

	classes := make([]Class, 0)

	for _, valuePair := range parseColonSeparatedLines(lines) {
		class := Class{
			Name:        valuePair[0],
			Description: valuePair[1],
		}

		classes = append(classes, class)
	}

	return classes, nil
}

// ListClass returns all devices in the specified device setup classes. Valid
// on local and remote computers.
//
// Notes
// Each entry in a setup class display represents one device. The entry consists
// of the unique instance name and a description of the device.
//
// To find the setup class of a particular device, use the Stack() function.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-listclass for more information.
func (dc *DevCon) ListClass() (map[string][]Device, error) {
	lines, err := dc.run(commandListClass)
	if err != nil {
		return nil, err
	}

	return parseListClass(lines), nil
}

// ClassFilter adds, deletes, displays, and changes the order of filter drivers
// for a device setup class. Valid only on the local computer.
//
// The class specifies the device setup class. The filter can be upper to indicate
// that the specified drivers are upper-class filter drivers, or lower to
// indicate that the specified drivers are lower-class filter drivers.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-classfilter for more information.
func (dc *DevCon) ClassFilter(class string, filter ClassFilter, drivers ...string) error {
	// TODO: Change params to allow for flexible input.
	args := []string{class, string(filter)}
	args = append(args, drivers...)

	lines, err := dc.run(commandClassFilter, args...)
	if err != nil {
		return err
	}

	// TODO: Parse lines.
	dc.logResults(lines)

	return nil
}

func parseListClass(lines []string) map[string][]Device {
	devicesByClassMap := make(map[string][]Device)

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
					class = matches[1]
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

		devicesByClassMap[class] = devices
	}

	return devicesByClassMap
}
