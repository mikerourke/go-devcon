package devcon

import "strings"

// HwID contains the hardware IDs and compatible IDs for a device.
type HwID struct {
	Device Device `json:"device"`

	// HardwareIDs is a vendor-defined identification string that Windows uses to match a device to an INF file.
	// In most cases, a device has more than one hardware ID associated with it.
	// Typically, a list of hardware IDs is sorted from most to least suitable for a device.
	//
	// See https://docs.microsoft.com/en-us/windows-hardware/drivers/install/hardware-ids for more information.
	HardwareIDs []string `json:"hardwareIds"`

	// CompatibleIDs are the vendor-defined identification strings that Windows uses to match a device to an INF file.
	//
	// See https://docs.microsoft.com/en-us/windows-hardware/drivers/install/compatible-ids for more information.
	CompatibleIDs []string `json:"compatibleIds"`
}

// HwIDs returns HwID records containing the hardware IDs, compatible IDs, and
// device instance IDs of the specified devices. Valid on local and remote
// computers.
//
// Example
//	dc := devcon.New("path\to\devcon.exe")
//	dc.OnRemote("server01").HwIDs("acpi*", "=usb")
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-hwids for more information.
func (dc *DevCon) HwIDs(ids ...string) ([]HwID, error) {
	lines, err := dc.run(commandHwIDs, ids...)
	if err != nil {
		return nil, err
	}

	return parseHwIDs(lines), nil
}

// SetHwID adds, deletes, and changes the order of hardware IDs of root-enumerated
// devices on a local or remote computer.
//
// Notes
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
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-sethwid for more information.
func (dc *DevCon) SetHwID(idsOrClasses []string, hardwareIds []string) error {
	args := idsOrClasses
	args = append(args, ":=")
	for _, hardwareId := range hardwareIds {
		args = append(args, hardwareId)
	}

	lines, err := dc.run(commandSetHwID, args...)
	if err != nil {
		return err
	}

	// TODO: Parse lines.
	dc.printResults(lines)

	return nil
}

func parseHwIDs(lines []string) []HwID {
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
			HardwareIDs:   make([]string, 0),
			CompatibleIDs: make([]string, 0),
		}
		search := None

		for lineIndex := groupStart; lineIndex < groupEnd; lineIndex++ {
			thisLine := lines[lineIndex]

			if lineIndex == groupStart {
				hwid.Device.ID = thisLine

				search = None
			} else if lineIndex == groupStart+1 {
				nameParams := parseParams(reName, thisLine)

				if name, ok := nameParams["Name"]; ok {
					hwid.Device.Name = name
				}
			} else if strings.Contains(thisLine, "Hardware ID") {
				search = HW
			} else if strings.Contains(thisLine, "Compatible ID") {
				search = Compat
			} else {
				idLine := strings.Trim(thisLine, " ")

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

	return hwids
}
