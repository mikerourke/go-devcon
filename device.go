package devcon

// Device represents a device attached to the computer.
type Device struct {
	// ID is a string reported by a device’s enumerator. A device has only one
	// device ID. A device ID has the same format as a hardware ID.
	//
	// // See https://docs.microsoft.com/en-us/windows-hardware/drivers/install/device-ids for more information.
	ID string `json:"id"`

	// Name is the name of the device.
	Name string `json:"name"`
}

// Find returns all devices that are currently attached to the computer. Includes
// the device instance ID and device name. Valid on local and remote computers.
//
// Notes
// You can use Find to find devices that are not currently attached to the computer
// by specifying the full device instance ID of the device instead of a hardware
// ID or ID pattern. Specifying the full device instance ID overrides the
// restriction on the Find operation that limits it to attached devices.
//
// Calling Find with a single class argument returns the same results as the
// ListClass() function.
//
// To find all devices, including those that are not currently attached to the
// computer, use the FindAll() function.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-find for more information.
func (dc *DevCon) Find(matchers ...string) ([]Device, error) {
	lines, err := dc.run(commandFind, matchers...)
	if err != nil {
		return nil, err
	}

	devices := make([]Device, 0)

	for _, valuePair := range parseColonSeparatedLines(lines) {
		device := Device{
			ID:   valuePair[0],
			Name: valuePair[1],
		}

		devices = append(devices, device)
	}

	return devices, nil
}

// FindAll returns all devices on the computer, including devices that were once
// attached to the computer but have been detached or moved. (These are known as
// non-present devices or phantom devices.) The FindAll operation also returns devices
// that are enumerated differently as a result of a BIOS change. Valid on local
// and remote computers.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-findall for more information.
func (dc *DevCon) FindAll() ([]Device, error) {
	lines, err := dc.run(commandFindAll)
	if err != nil {
		return nil, err
	}

	devices := make([]Device, 0)

	for _, valuePair := range parseColonSeparatedLines(lines) {
		device := Device{
			ID:   valuePair[0],
			Name: valuePair[1],
		}

		devices = append(devices, device)
	}

	return devices, nil
}

// Disable disables devices on the computer. Valid only on the local computer.
//
// To disable a device means that the device remains physically connected to the
// computer, but its driver is unloaded from memory and its resources are freed
// so that the device cannot be used.
//
// Notes
// This will disable the device even if the device is already disabled. Before
// and after disabling a device, use Status() to verify the device status.
//
// Before using an ID pattern to disable a device, determine which devices will
// be affected. To do so, pass the pattern to the Status() function:
//	dc.Status("USB\*")
// Or with the HwIDs() function:
// 	dc.HwIDs("USB\*")
//
// The system might need to be rebooted to make this change effective. To
// reboot the system, add ConditionalReboot() before Disable().
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-disable for more information.
func (dc *DevCon) Disable(matchers ...string) error {
	lines, err := dc.run(commandDisable, matchers...)
	if err != nil {
		return err
	}

	dc.logResults(lines)

	return err
}

// Enable enables devices on the computer. Valid only on the local computer.
//
// To enable a device means that the device driver is loaded into memory and the
// device is ready for use.
//
// Notes
// This will enable the device even if it is already enabled. Before and after
// enabling a device, use Status() to verify the device status.
//
// The system might need to be rebooted to make this change effective. To
// reboot the system, add ConditionalReboot() before Enable().
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-enable for more information.
func (dc *DevCon) Enable(matchers ...string) error {
	lines, err := dc.run(commandEnable, matchers...)
	if err != nil {
		return err
	}

	// TODO: Parse lines.
	dc.logResults(lines)

	return err
}

// Install creates a new, root-enumerated devnode for a non-Plug and Play device
// and installs its supporting software. Valid only on the local computer.
// The infFilePath parameter is the fully qualified path and
// name of the INF file for the driver package.
//
// The hardwareID specifies a hardware ID for the device. It must exactly match
// the hardware ID of the device. Patterns are not valid. Do not use a single
// quote character (') to indicate a literal value. For more information, see
// https://docs.microsoft.com/en-us/windows-hardware/drivers/install/hardware-ids
// and https://docs.microsoft.com/en-us/windows-hardware/drivers/install/device-identification-strings.
//
// Notes
// The system might need to be rebooted to make this change effective. To reboot
// the system, add ConditionalReboot() to Install().
//
// You cannot use this function for Plug and Play devices.
//
// This operation creates a new non-Plug and Play device node. Then, it uses
// Update to install drivers for the newly added device.
//
// If any step of the installation operation fails, returns an error and does
// not proceed with the driver installation.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-install for more information.
func (dc *DevCon) Install(infFilePath string, hardwareID string) error {
	lines, err := dc.run(commandInstall, infFilePath, hardwareID)
	if err != nil {
		return err
	}

	// TODO: Parse lines.
	dc.logResults(lines)

	return nil
}

// Remove removes the device from the device tree and deletes the device stack
// for the device. As a result of these actions, child devices are removed from
// the device tree and the drivers that support the device are unloaded.
//
// This operation does not delete the device driver or any files installed for
// the device. After removing the device from the device tree, the files remain
// and the device is still represented internally as a non-present device that
// can be re-enumerated.
//
// Valid only on the local computer.
//
// Notes
// The system might need to be rebooted to make this change effective. To reboot
// the system, add ConditionalReboot() before Remove().
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-remove for more information.
func (dc *DevCon) Remove(matchers ...string) error {
	lines, err := dc.run(commandRemove, matchers...)
	if err != nil {
		return err
	}

	// TODO: Parse lines.
	dc.logResults(lines)

	return nil
}

// Rescan uses Windows Plug and Play features to update the device list for the
// computer. Valid on local and remote computers.
//
// Notes
// Rescanning can cause the Plug and Play manager to detect new devices and to
// install device drivers without warning.
//
// Rescanning can detect some non-Plug and Play devices, particularly those
// that cannot notify the system when they are installed, such as parallel-port
// devices and serial-port devices. As a result, you must have Administrator
// privileges to call this function.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-rescan for more information.
func (dc *DevCon) Rescan() error {
	lines, err := dc.run(commandRescan)
	if err != nil {
		return err
	}

	// TODO: Parse lines.
	dc.logResults(lines)

	return nil
}

// Restart stops and restarts the specified devices. Valid only on the local
// computer.
//
// Notes
// The system might need to be rebooted to make this change effective. To
// reboot the system, add ConditionalReboot() before Restart().
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-restart for more information.
func (dc *DevCon) Restart(matches ...string) error {
	lines, err := dc.run(commandRestart, matches...)
	if err != nil {
		return err
	}

	// TODO: Parse lines.
	dc.logResults(lines)

	return nil
}
