package devcon

import "strings"

// HwID contains the hardware IDs and compatible IDs for a device.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/install/hardware-ids for more information
// about Hardware IDs.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/install/compatible-ids for more information
// about Compatible IDs.
type HwID struct {
	Device Device `json:"device"`

	// HardwareIDs is a vendor-defined identification string that Windows uses
	// to match a device to an INF file. In most cases, a device has more than
	// one hardware ID associated with it. Typically, a list of hardware IDs is
	// sorted from most to least suitable for a device.
	HardwareIDs []string `json:"hardwareIds"`

	// CompatibleIDs are the vendor-defined identification strings that Windows
	// uses to match a device to an INF file.
	CompatibleIDs []string `json:"compatibleIds"`
}

// HwIDs returns HwID records containing the hardware IDs, compatible IDs, and
// device instance IDs of the specified devices. Valid on local and remote
// computers.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-hwids for more information.
func (dc *DevCon) HwIDs(idsOrClasses ...string) ([]HwID, error) {
	lines, err := dc.run(commandHwIDs, idsOrClasses...)
	if err != nil {
		return nil, err
	}

	type searchStatus int

	const (
		None searchStatus = iota
		HW
		Compat
	)

	groupIndices := make([]int, 0)

	for index, line := range lines {
		if !strings.HasPrefix(line, " ") {
			groupIndices = append(groupIndices, index)
		}
	}

	groupIndices = append(groupIndices, len(lines))

	hwids := make([]HwID, 0)

	for index, groupStart := range groupIndices {
		nextIndex := index + 1
		if len(groupIndices) == nextIndex {
			break
		}

		groupEnd := groupIndices[nextIndex]

		hwid := HwID{
			Device: Device{
				ID:   "",
				Name: "",
			},
			HardwareIDs:   make([]string, 0),
			CompatibleIDs: make([]string, 0),
		}
		search := None

		for lineIndex := groupStart; lineIndex < groupEnd; lineIndex++ {
			line := lines[lineIndex]

			switch {
			case lineIndex == groupStart:
				hwid.Device.ID = line

				search = None

			case lineIndex == groupStart+1:
				nameParams := parseParams(reName, line)

				if name, ok := nameParams["Name"]; ok {
					hwid.Device.Name = name
				}

			case strings.Contains(line, "Hardware ID"):
				search = HW

			case strings.Contains(line, "Compatible ID"):
				search = Compat

			default:
				idLine := trimSpaces(line)

				if search == HW {
					hwid.HardwareIDs = append(hwid.HardwareIDs, idLine)
				} else if search == Compat {
					hwid.CompatibleIDs = append(hwid.CompatibleIDs, idLine)
				}
			}
		}

		if hwid.Device.Name != "" {
			hwids = append(hwids, hwid)
		}
	}

	return hwids, nil
}

// SetHwID adds, deletes, and changes the order of hardware IDs of root-enumerated
// devices.
//
// A root-enumerated device is a device that appears in the ROOT registry
// subkey (HKEY_LOCAL_MACHINE\System\ControlSet\Enum\ROOT).
//
// You can specify multiple hardware IDs in each command. The ! (delete) parameter
// applies only to the hardware ID that it prefixes. The other symbol parameters
// apply to all hardware IDs that follow until the next symbol parameter in the
// command.
//
// SetHwID moves, rather than adds, a hardware ID if the specified hardware ID
// already exists in the list of hardware IDs for the device.
//
// Running with the WithRemoteComputer() option is allowed.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-sethwid for more information.
func (dc *DevCon) SetHwID(idsOrClasses []string, hardwareIds []string) error {
	args := idsOrClasses
	args = append(args, ":=")
	args = append(args, hardwareIds...)

	lines, err := dc.run(commandSetHwID, args...)
	if err != nil {
		return err
	}

	// TODO: Parse lines.
	dc.logResults(lines)

	return nil
}
