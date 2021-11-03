package devcon

import (
	"errors"
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

// DevCon wraps the Device Console package and `devcon.exe` executable and
// provides commands for interacting with the utility and performing operations.
type DevCon struct {
	// ExeFilePath is the path to `devcon.exe` on the computer.
	ExeFilePath string

	isReboot   bool
	remotePath string
}

// New returns a new instance of DevCon that can be used to run commands and
// queries.
func New(exeFilePath string) *DevCon {
	return &DevCon{
		ExeFilePath: exeFilePath,
	}
}

// ConditionalReboot should be set on the DevCon instance if the computer should
// be rebooted after running the command. If specified for a command that doesn't
// allow a conditional reboot, the command will not be run.
func (dc *DevCon) ConditionalReboot() *DevCon {
	dc.isReboot = true

	return dc
}

// OnRemoteComputer should be set on the DevCon instance if the command should
// be run on a remote computer. If specified for a command that cannot be
// run on a remote computer, the command will not be run.
//
// Important
// Ensure there are no leading backslashes specified in the remote path.
func (dc *DevCon) OnRemoteComputer(remotePath string) *DevCon {
	dc.remotePath = remotePath

	return dc
}

// run executes the `devcon.exe` tool with the specified command and args.
func (dc *DevCon) run(command command, args ...string) ([]string, error) {
	if dc.isReboot && !command.CanReboot() {
		return nil, fmt.Errorf(
			"the %s command does not allow a conditional reboot, remove .ConditionalReboot() to proceed",
			command)
	}

	if dc.remotePath != "" && !command.CanBeRemote() {
		return nil, fmt.Errorf(
			"the %s command cannot be ran on a remote computer, remove .OnRemoteComputer() to proceed",
			command)
	}

	allArgs := make([]string, 0)

	if dc.remotePath != "" {
		if !strings.HasPrefix(dc.remotePath, `\`) {
			return nil, errors.New("the remote computer name must have leading backslashes")
		}

		allArgs = append(allArgs, fmt.Sprintf(`/m:%s`, dc.remotePath))
	}

	if dc.isReboot {
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
	dc.isReboot = false
	dc.remotePath = ""

	// Read and parse the contents of the file associated with the command
	// from the `testdata` directory.
	if dc.ExeFilePath == "" {
		fmt.Println("No path specified for devcon.exe, using test data")

		return readTestDataFile(command)
	}

	out, err := exec.Command(dc.ExeFilePath, allArgs...).Output()
	if err != nil {
		return nil, fmt.Errorf("error running %s command: %w", command, err)
	}

	lines := splitLines(string(out))

	return lines, nil
}

// logResults logs out the results of a command to the console.
func (dc *DevCon) logResults(lines []string) {
	for _, line := range lines {
		fmt.Println(line)
	}
}

func readTestDataFile(command command) ([]string, error) {
	fileName := fmt.Sprintf("%s.txt", command.String())

	path := filepath.Join("testdata", fileName)

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, fmt.Errorf("error reading test data file: %w", err)
	}

	return splitLines(string(bytes)), nil
}
