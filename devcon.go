// Package devcon wraps the Windows Device Console (devcon.exe) utility.
//
// The Windows Device Console is a command-line tool that displays detailed
// information about devices on computers running Windows. You can use DevCon
// to enable, disable, install, configure, and remove devices.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon for more information.
package devcon

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// DevCon is the main entry point for the utility. It is created with New().
type DevCon struct {
	// ExeFilePath is the path to `devcon.exe` on the computer.
	ExeFilePath string

	// IsRebooted indicates that the computer will be conditionally rebooted
	// after running a command.
	IsRebooted bool

	// RemotePath is the path to a remote computer.
	RemotePath string
}

// New returns a new instance of DevCon that can be used to run commands and
// queries.
//
// The path to `devcon.exe` must be specified because this package does not
// ship with a copy of the executable and versions may vary by Windows version.
func New(exeFilePath string) *DevCon {
	return &DevCon{
		ExeFilePath: exeFilePath,
		IsRebooted:  false,
		RemotePath:  "",
	}
}

// WithConditionalReboot should be set on the DevCon instance if the computer
// should be rebooted after running the command. If specified for a command that
// doesn't allow a conditional reboot, the command will not be run.
//
// Note that the computer will only be rebooted if a reboot is required.
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-general-commands for more information.
//
// # Usage
//
// Add to the DevCon instance before the command.
//
//	dc.WithConditionalReboot().Enable("*")
func (dc *DevCon) WithConditionalReboot() *DevCon {
	dc.IsRebooted = true

	return dc
}

// WithRemoteComputer should be set on the DevCon instance if the command should
// be run on a remote computer. If specified for a command that cannot be
// run on a remote computer, the command will not be run.
//
// To call a method for a remote computer, the Group Policy setting must allow
// the Plug and Play service to run on the remote computer. On computers that
// run Windows Vista and Windows 7, the Group Policy disables remote access to
// the service by default. On computers that run WDK 8.1 and WDK 8, the remote
// access is unavailable.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-general-commands for more information.
//
// # Usage
//
// Add to the DevCon instance before the command.
//
//	dc.WithRemoteComputer(`\\server01`).Enable("USB\VID_0403&PID_6001\AB0K0SRH")
func (dc *DevCon) WithRemoteComputer(remotePath string) *DevCon {
	dc.RemotePath = remotePath

	return dc
}

// run executes the `devcon.exe` tool with the specified command and args.
func (dc *DevCon) run(command command, args ...string) ([]string, error) {
	if dc.IsRebooted && !command.CanReboot() {
		return nil, fmt.Errorf(
			"the %s command does not allow a conditional reboot, remove .WithConditionalReboot() to proceed",
			command)
	}

	if dc.RemotePath != "" && !command.CanBeRemote() {
		return nil, fmt.Errorf(
			"the %s command cannot be ran on a remote computer, remove .WithRemoteComputer() to proceed",
			command)
	}

	allArgs := make([]string, 0)

	if dc.RemotePath != "" {
		if !strings.HasPrefix(dc.RemotePath, `\`) {
			return nil, errors.New("the remote computer name must have leading backslashes")
		}

		allArgs = append(allArgs, `/m:`+dc.RemotePath)
	}

	if dc.IsRebooted {
		allArgs = append(allArgs, "/r")
	}

	allArgs = append(allArgs, command.String())

	for _, arg := range args {
		if strings.Contains(arg, "&") {
			allArgs = append(allArgs, fmt.Sprintf(`"%s"`, arg))
		} else {
			allArgs = append(allArgs, arg)
		}
	}

	// Reset these to their defaults to ensure they don't get applied to any
	// subsequent commands.
	dc.IsRebooted = false
	dc.RemotePath = ""

	if dc.ExeFilePath == "" {
		return nil, errors.New("invalid devcon.exe path specified")
	}

	if _, err := os.Stat(dc.ExeFilePath); os.IsNotExist(err) {
		return nil, fmt.Errorf("%s does not exist", dc.ExeFilePath)
	}

	output, err := exec.Command(dc.ExeFilePath, allArgs...).Output()

	lines := splitLines(string(output))

	if err != nil {
		return lines, fmt.Errorf("error running %s command: %s", command, err)
	}

	return lines, nil
}

// substrInLines returns the index of the line in the specified lines slice
// where the matching substr was found. If the substr was not found, return -1.
func substrInLines(lines []string, substr string) int {
	for index, line := range lines {
		if strings.Contains(line, substr) {
			return index
		}
	}

	return -1
}

// concatIdsOrClasses combines the first ID or class (required for a command)
// with any additional variadic inputs to the method.
func concatIdsOrClasses(idOrClass string, idsOrClasses ...string) []string {
	all := []string{idOrClass}

	if len(idsOrClasses) != 0 {
		all = append(all, idsOrClasses...)
	}

	return all
}

// logResults logs out the results of a command to the console.
func (dc *DevCon) logResults(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}
