package devcon

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
	"strings"
)

type Options struct {
	ID     string
	Reboot bool
}

type DevCon struct {
	ExeFilePath string

	isReboot bool
	remote   string
}

func New(exeFilePath string) *DevCon {
	return &DevCon{
		ExeFilePath: exeFilePath,
	}
}

func (dc *DevCon) ConditionalReboot() *DevCon {
	dc.isReboot = true

	return dc
}

func (dc *DevCon) OnRemoteComputer(remote string) *DevCon {
	dc.remote = remote

	return dc
}

func (dc *DevCon) run(command command, args ...string) ([]string, error) {
	if dc.isReboot && !command.CanReboot() {
		return nil, fmt.Errorf("the %s command does not allow a conditional reboot, remove .ConditionalReboot() to proceed", command)
	}

	if dc.remote != "" && !command.CanBeRemote() {
		return nil, fmt.Errorf("the %s command cannot be ran on a remote computer, remove .OnRemoteComputer() to proceed", command)
	}

	// return readTestDataFile(command)

	allFlags := make([]string, 0)

	if dc.remote != "" {
		if strings.HasPrefix(dc.remote, `\`) {
			return nil, fmt.Errorf(`the remote computer name cannot have leading backslashes`)
		}

		allFlags = append(allFlags, fmt.Sprintf(`/m:\\%s`, dc.remote))
	}

	if dc.isReboot {
		allFlags = append(allFlags, "/r")
	}

	allFlags = append(allFlags, command.String())

	allFlags = append(allFlags, args...)

	return readTestDataFile(command)

	out, err := exec.Command(dc.ExeFilePath, allFlags...).Output()

	dc.isReboot = false
	dc.remote = ""

	if err != nil {
		return nil, err
	}

	lines := splitLines(string(out))

	return lines, err
}

func readTestDataFile(command command) ([]string, error) {
	fileName := fmt.Sprintf("%s.txt", command.String())

	path := filepath.Join("testdata", fileName)

	bytes, err := ioutil.ReadFile(path)
	if err != nil {
		return nil, err
	}

	return splitLines(string(bytes)), nil
}
