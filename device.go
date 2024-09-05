package devcon

import (
	"errors"
	"fmt"
	"strings"
)

// Device represents a device attached to the computer.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/install/device-ids for more information about Device IDs.
type Device struct {
	// ID is a string reported by a device’s enumerator. A device has only one
	// device ID. A device ID has the same format as a hardware ID.
	ID string `json:"id"`

	// Name is the name of the device.
	Name string `json:"name"`
}

// DeviceRestartStatus indicates the restart status of a device when Restart()
// is called.
type DeviceRestartStatus struct {
	// ID is the ID of the device being restarted.
	ID string `json:"id"`

	// WasRestarted indicates if the restart operation was successful.
	WasRestarted bool
}

// Disable disables devices on the computer.
//
// To disable a device means that the device remains physically connected to the
// computer, but its driver is unloaded from memory and its resources are freed
// so that the device cannot be used.
//
// This will disable the device even if the device is already disabled. Before
// and after disabling a device, use Status() to verify the device status.
//
// Before using an ID pattern to disable a device, determine which devices will
// be affected. To do so, pass the pattern to the Status() function:
//
//	dc.Status("USB\*")
//
// Or with the HwIDs() function:
//
//	dc.HwIDs("USB\*")
//
// The system might need to be rebooted to make this change effective. To reboot
// the system if required, use:
//
//	dc.WithConditionalReboot().Disable()
//
// Cannot be run with the WithRemoteComputer() option.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-disable for more information.
func (dc *DevCon) Disable(idOrClass string, idsOrClasses ...string) error {
	allIdsOrClasses := concatIdsOrClasses(idOrClass, idsOrClasses...)

	lines, err := dc.run(commandDisable, allIdsOrClasses...)
	if err != nil {
		return err
	}

	dc.logResults(lines)

	return err
}

// Enable enables devices on the computer.
//
// To enable a device means that the device driver is loaded into memory and the
// device is ready for use.
//
// This will enable the device even if it is already enabled. Before and after
// enabling a device, use Status() to verify the device status.
//
// The system might need to be rebooted to make this change effective. To reboot
// the system if required, use:
//
//	dc.WithConditionalReboot().Enable()
//
// Cannot be run with the WithRemoteComputer() option.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-enable for more information.
func (dc *DevCon) Enable(idOrClass string, idsOrClasses ...string) error {
	allIdsOrClasses := concatIdsOrClasses(idOrClass, idsOrClasses...)

	lines, err := dc.run(commandEnable, allIdsOrClasses...)
	if err != nil {
		return err
	}

	// TODO: Parse lines.
	dc.logResults(lines)

	return err
}

// Find returns a slice of Device instances that are currently attached to the computer.
//
// You can use Find to find devices that are not currently attached to the computer
// by specifying the full device instance ID of the device instead of a hardware
// ID or ID pattern. Specifying the full device instance ID overrides the
// restriction on the Find operation that limits it to attached devices.
//
// Calling Find with a single class argument returns the same results as the
// ListClass() method.
//
// To find all devices, including those that are not currently attached to the
// computer, use the FindAll() function.
//
// Running with the WithRemoteComputer() option is allowed.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-find for more information.
func (dc *DevCon) Find(idOrClass string, idsOrClasses ...string) ([]Device, error) {
	allIdsOrClasses := concatIdsOrClasses(idOrClass, idsOrClasses...)

	lines, err := dc.run(commandFind, allIdsOrClasses...)
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
// non-present devices or phantom devices.).
//
// It also returns devices that are enumerated differently as a result of a
// BIOS change.
//
// Running with the WithRemoteComputer() option is allowed.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-findall for more information.
func (dc *DevCon) FindAll(idOrClass string, idsOrClasses ...string) ([]Device, error) {
	allIdsOrClasses := concatIdsOrClasses(idOrClass, idsOrClasses...)

	lines, err := dc.run(commandFindAll, allIdsOrClasses...)
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

// Install creates a new, root-enumerated devnode for a non-Plug and Play device
// and installs its supporting software.
//
// The infFilePath parameter is the fully qualified path and name of the INF file
// for the driver package.
//
// The hardwareID specifies a hardware ID for the device. It must exactly match
// the hardware ID of the device. Patterns are not valid. Do not use a single
// quote character (') to indicate a literal value. For more information, see
// https://docs.microsoft.com/en-us/windows-hardware/drivers/install/hardware-ids
// and https://docs.microsoft.com/en-us/windows-hardware/drivers/install/device-identification-strings.
//
// You cannot use this function for Plug and Play devices. This operation creates
// a new non-Plug and Play device node. Then, it uses Update to install drivers
// for the newly added device. If any step of the installation operation fails,
// returns an error and does not proceed with the driver installation.
//
// The system might need to be rebooted to make this change effective. To reboot
// the system if required, use:
//
//	dc.WithConditionalReboot().Install()
//
// Cannot be run with the WithRemoteComputer() option.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-install for more information.
func (dc *DevCon) Install(infFilePath string, hardwareID string) error {
	lines, err := dc.run(commandInstall, infFilePath, hardwareID)
	if err != nil {
		return err
	}

	if substrInLines(lines, "success") != -1 {
		return errors.New("install command was not successful")
	}

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
// The system might need to be rebooted to make this change effective. To reboot
// the system if required, use:
//
//	dc.WithConditionalReboot().Remove()
//
// Cannot be run with the WithRemoteComputer() option.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-remove for more information.
func (dc *DevCon) Remove(idsOrClasses ...string) error {
	lines, err := dc.run(commandRemove, idsOrClasses...)
	if err != nil {
		return err
	}

	// TODO: Parse lines.
	dc.logResults(lines)

	return nil
}

// Rescan uses Windows Plug and Play features to update the device list for the
// computer.
//
// Rescanning can cause the Plug and Play manager to detect new devices and to
// install device drivers without warning.
//
// Rescanning can detect some non-Plug and Play devices, particularly those
// that cannot notify the system when they are installed, such as parallel-port
// devices and serial-port devices. As a result, you must have Administrator
// privileges to call this function.
//
// The system might need to be rebooted to make this change effective. To reboot
// the system if required, use:
//
//	dc.WithConditionalReboot().Rescan()
//
// Running with the WithRemoteComputer() option is allowed.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-rescan for more information.
func (dc *DevCon) Rescan() error {
	lines, err := dc.run(commandRescan)
	if err != nil {
		return err
	}

	if substrInLines(lines, "completed") == -1 {
		return fmt.Errorf("unable to rescan: %s", strings.Join(lines, ". "))
	}

	return nil
}

// Restart stops and restarts the specified devices.
//
// The system might need to be rebooted to make this change effective. To reboot
// the system if required, use:
//
//	dc.WithConditionalReboot().Restart()
//
// Cannot be run with the WithRemoteComputer() option.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-restart for more information.
func (dc *DevCon) Restart(idOrClass string, idsOrClasses ...string) ([]DeviceRestartStatus, error) {
	allIdsOrClasses := concatIdsOrClasses(idOrClass, idsOrClasses...)

	lines, err := dc.run(commandRestart, allIdsOrClasses...)
	if err != nil {
		return nil, err
	}

	statuses := make([]DeviceRestartStatus, 0)

	valuePairs := parseColonSeparatedLines(lines)
	for _, valuePair := range valuePairs {
		statuses = append(statuses, DeviceRestartStatus{
			ID:           valuePair[0],
			WasRestarted: valuePair[1] == "Restarted",
		})
	}

	return statuses, nil
}
