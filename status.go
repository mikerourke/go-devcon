package devcon

import "strings"

// Status indicates the current status of a device driver.
type Status string

const (
	// StatusRunning indicates that the device driver is running.
	StatusRunning Status = "running"

	// StatusStopped indicates that the device driver is stopped.
	StatusStopped Status = "stopped"

	// StatusDisabled indicates that the device driver is disabled.
	StatusDisabled Status = "disabled"

	// StatusUnknown indicates that the status could not be queried or is
	// unknown.
	StatusUnknown Status = "unknown"
)

// DriverStatus contains details of the status of a device driver.
type DriverStatus struct {
	// Device is the device that corresponds to the driver.
	Device Device `json:"device"`

	// Status is the current status of the device driver.
	Status Status `json:"status"`
}

// Status returns the status (running, stopped, or disabled) of the driver for
// devices on the computer.
//
// If the status of the device cannot be determined, such as when the device is
// no longer attached to the computer, the status from the status display.
//
// Running with the WithRemoteComputer() option is allowed.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-status for more information.
func (dc *DevCon) Status() ([]DriverStatus, error) {
	lines, err := dc.run(commandStatus)
	if err != nil {
		return nil, err
	}

	return parseStatus(lines), nil
}

func parseStatus(lines []string) []DriverStatus {
	groupIndices := make([]int, 0)

	for index, line := range lines {
		if !strings.HasPrefix(line, " ") {
			groupIndices = append(groupIndices, index)
		}
	}

	groupIndices = append(groupIndices, len(lines))

	statuses := make([]DriverStatus, 0)

	for index, groupStart := range groupIndices {
		nextIndex := index + 1
		if len(groupIndices) == nextIndex {
			break
		}

		groupEnd := groupIndices[nextIndex]

		status := DriverStatus{}

		for lineIndex := groupStart; lineIndex < groupEnd; lineIndex++ {
			line := lines[lineIndex]

			switch {
			case lineIndex == groupStart:
				status.Device.ID = line

			case lineIndex == groupStart+1:
				nameParams := parseParams(reName, line)

				if name, ok := nameParams["Name"]; ok {
					status.Device.Name = name
				}

			default:
				statusLine := trimSpaces(line)

				switch {
				case strings.Contains(statusLine, "running"):
					status.Status = StatusRunning

				case strings.Contains(statusLine, "stopped"):
					status.Status = StatusStopped

				case strings.Contains(statusLine, "disabled"):
					status.Status = StatusDisabled

				default:
					status.Status = StatusUnknown
				}
			}
		}

		if status.Device.Name != "" {
			statuses = append(statuses, status)
		}
	}

	return statuses
}
