package devcon

import (
	"regexp"
	"strings"
)

var (
	reSetupClass = regexp.MustCompile(`Setup Class: {(?P<GUID>.*)} (?P<Name>.*)`)
)

// DriverStack contains driver stack details for a device.
type DriverStack struct {
	// Device is the device that the stack is associated with.
	Device Device `json:"device"`

	// SetupClassGUID is the GUID of the setup class.
	SetupClassGUID string `json:"setupClassGuid"`

	// SetupClassName is the name of the setup class.
	SetupClassName string `json:"setupClassName"`

	// ControllingService is the Windows Service that controls the driver.
	ControllingService string `json:"controllingService"`

	// UpperFilters represent a configuration that monitors all IRP traffic
	// into or out of the device driver within the device's driver stack,
	// regardless of whether the driver processed the IRP or passed it through
	// to lower device drivers.
	UpperFilters string `json:"upperFilters"`

	// LowerFilters represent a configuration that monitors all IRP traffic into
	// or out of the device driver from lower drivers within the device's driver stack.
	LowerFilters string `json:"lowerFilters"`
}

// Stack returns the expected driver stack for the specified devices, and the
// GUID and the name of the device setup class for each device. Although the
// actual driver stack typically matches the expected stack, variations are
// possible.
//
// To investigate a device problem, compare the expected driver stack from the
// stack operation with the actual drivers that the device uses, as returned
// from DriverFiles().
//
// Running with the WithRemoteComputer() option is allowed.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-stack for more information.
func (dc *DevCon) Stack() ([]DriverStack, error) {
	lines, err := dc.run(commandStack)
	if err != nil {
		return nil, err
	}

	type searchStatus int

	const (
		None searchStatus = iota
		UpperFilter
		Service
		LowerFilter
	)

	groupIndices := make([]int, 0)

	for index, line := range lines {
		if !strings.HasPrefix(line, " ") {
			groupIndices = append(groupIndices, index)
		}
	}

	groupIndices = append(groupIndices, len(lines))

	stacks := make([]DriverStack, 0)

	for index, groupStart := range groupIndices {
		nextIndex := index + 1
		if len(groupIndices) == nextIndex {
			break
		}

		groupEnd := groupIndices[nextIndex]

		stack := DriverStack{
			Device: Device{
				ID:   "",
				Name: "",
			},
			SetupClassGUID:     "",
			SetupClassName:     "",
			ControllingService: "",
			UpperFilters:       "",
			LowerFilters:       "",
		}

		search := None

		for lineIndex := groupStart; lineIndex < groupEnd; lineIndex++ {
			line := trimSpaces(lines[lineIndex])

			switch {
			case lineIndex == groupStart:
				stack.Device.ID = line

			case lineIndex == groupStart+1:
				nameParams := parseParams(reName, line)

				if name, ok := nameParams["Name"]; ok {
					stack.Device.Name = name
				}

			case strings.Contains(line, "Setup Class"):
				params := parseParams(reSetupClass, line)

				if guid, ok := params["GUID"]; ok {
					stack.SetupClassGUID = guid
				}

				if name, ok := params["Name"]; ok {
					stack.SetupClassName = name
				}

			case strings.Contains(line, "pper filters:"):
				search = UpperFilter

			case strings.Contains(line, "service:"):
				search = Service

			case strings.Contains(line, "ower filters:"):
				search = LowerFilter

			default:
				search = None
			}

			//goland:noinspection GoDfaConstantCondition
			switch search {
			case UpperFilter:
				stack.UpperFilters = line

			case Service:
				stack.ControllingService = line

			case LowerFilter:
				stack.LowerFilters = line

			case None: // Do nothing
			}
		}

		if stack.Device.Name != "" {
			stacks = append(stacks, stack)
		}
	}

	return stacks, nil
}
