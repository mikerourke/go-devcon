package devcon

import (
	"fmt"
	"io/ioutil"
	"os/exec"
	"path/filepath"
)

type command string

const (
	commandClasses     command = "classes"
	commandClassFilter command = "classfilter"
	commandDisable     command = "disable"
	commandDpAdd       command = "dp_add"
	commandDpDelete    command = "dp_delete"
	commandDpEnum      command = "dp_enum"
	commandDriverFiles command = "driverfiles"
	commandDriverNodes command = "drivernodes"
	commandEnable      command = "enable"
	commandFind        command = "find"
	commandFindAll     command = "findall"
	commandHwIDs       command = "hwids"
	commandInstall     command = "install"
	commandListClass   command = "listclass"
	commandReboot      command = "reboot"
	commandRemove      command = "remove"
	commandRescan      command = "rescan"
	commandResources   command = "resources"
	commandRestart     command = "restart"
	commandSetHwID     command = "sethwid"
	commandStack       command = "stack"
	commandStatus      command = "status"
	commandUpdate      command = "update"
	commandUpdateNI    command = "updateni"
)

func (c command) String() string {
	return string(c)
}

type DevCon struct {
	ExeFilePath string
}

func New(exeFilePath string) *DevCon {
	return &DevCon{
		ExeFilePath: exeFilePath,
	}
}

func (dc *DevCon) runWithoutArgs(command command) ([]string, error) {
	return readTestDataFile(command)

	out, err := exec.Command(dc.ExeFilePath, command.String()).Output()

	if err != nil {
		return nil, err
	}

	lines := splitLines(string(out))

	return lines, nil
}

func (dc *DevCon) runWithArgs(command command, args []string) ([]string, error) {
	return readTestDataFile(command)

	allFlags := []string{command.String()}
	allFlags = append(allFlags, args...)

	out, err := exec.Command(dc.ExeFilePath, allFlags...).Output()

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
