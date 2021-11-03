package devcon

import "strings"

type Status string

const (
	IsRunning  Status = "running"
	IsStopped  Status = "stopped"
	IsDisabled Status = "disabled"
	IsUnknown  Status = "unknown"
)

type DriverStatus struct {
	Device Device `json:"device"`
	Status Status `json:"status"`
}

// Status returns the status (running, stopped, or disabled) of the driver for
// devices on the computer. Valid on local and remote computers.
//
// Notes
// If the status of the device cannot be determined, such as when the device is
// no longer attached to the computer, the status from the status display.
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
					status.Status = IsRunning

				case strings.Contains(statusLine, "stopped"):
					status.Status = IsStopped

				case strings.Contains(statusLine, "disabled"):
					status.Status = IsDisabled

				default:
					status.Status = IsUnknown
				}
			}
		}

		if status.Device.Name != "" {
			statuses = append(statuses, status)
		}
	}

	return statuses
}
