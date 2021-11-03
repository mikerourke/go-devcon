package devcon

import (
	"regexp"
	"strings"
)

var reResource = regexp.MustCompile(`(?P<Name>.*)\s*: (?P<Value>.*)`)

// Resource represents an assignable and addressable bus paths, such as DMA
// channels, I/O ports, IRQ, and memory addresses.
type Resource struct {
	// Name is the name of the resource (e.g. IO or IRQ).
	Name string `json:"name"`

	// Value is the value of the resource (e.g. "ffa0-ffaf" or 14).
	Value string `json:"value"`
}

// DeviceResourceUsage describes the resources that a device is currently using.
type DeviceResourceUsage struct {
	// Device is the device associated with the resource usages.
	Device Device `json:"device"`

	// Resources are the resources currently being used by the device.
	Resources []Resource `json:"resources"`
}

// Resources returns the resources allocated to the specified devices. Resources
// are assignable and addressable bus paths, such as DMA channels, I/O ports,
// IRQ, and memory addresses. Valid on local and remote computers.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-resources for more information.
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

	groupIndices = append(groupIndices, len(lines))

	resourceUsages := make([]DeviceResourceUsage, 0)

	for index, groupStart := range groupIndices {
		nextIndex := index + 1
		if len(groupIndices) == nextIndex {
			break
		}

		groupEnd := groupIndices[nextIndex]

		resourceUsage := DeviceResourceUsage{
			Device:    Device{},
			Resources: make([]Resource, 0),
		}

		for lineIndex := groupStart; lineIndex < groupEnd; lineIndex++ {
			line := trimSpaces(lines[lineIndex])

			switch {
			case lineIndex == groupStart:
				resourceUsage.Device.ID = line

			case lineIndex == groupStart+1:
				nameParams := parseParams(reName, line)

				if name, ok := nameParams["Name"]; ok {
					resourceUsage.Device.Name = name
				}

			default:
				params := parseParams(reResource, line)

				resource := Resource{}
				if name, ok := params["Name"]; ok {
					resource.Name = trimSpaces(name)
				}

				if value, ok := params["Value"]; ok {
					resource.Value = trimSpaces(value)
				}

				if resource.Name != "" {
					resourceUsage.Resources = append(resourceUsage.Resources, resource)
				}
			}
		}

		if resourceUsage.Device.Name != "" {
			resourceUsages = append(resourceUsages, resourceUsage)
		}
	}

	return resourceUsages
}
