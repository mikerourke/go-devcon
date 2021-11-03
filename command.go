package devcon

type command string

const (
	// These commands represent queries, they don't make any changes to the
	// system.
	commandClasses     command = "classes"
	commandListClass   command = "listclass"
	commandDriverFiles command = "driverfiles"
	commandDriverNodes command = "drivernodes"
	commandHwIDs       command = "hwids"
	commandStack       command = "stack"
	commandStatus      command = "status"
	commandFind        command = "find"
	commandFindAll     command = "findall"

	// The commands represent actions, they _do_ make changes to the system.
	commandClassFilter command = "classfilter"
	commandDisable     command = "disable"
	commandDPAdd       command = "dp_add"
	commandDPDelete    command = "dp_delete"
	commandDPEnum      command = "dp_enum"
	commandEnable      command = "enable"
	commandInstall     command = "install"
	commandReboot      command = "reboot"
	commandRemove      command = "remove"
	commandRescan      command = "rescan"
	commandResources   command = "resources"
	commandRestart     command = "restart"
	commandSetHwID     command = "sethwid"
	commandUpdate      command = "update"
	commandUpdateNI    command = "updateni"
)

// String returns the string representation of the command (used to pass to the
// run operation).
func (c command) String() string {
	return string(c)
}

// CanReboot returns true if the associated command allows for conditional
// reboot.
func (c command) CanReboot() bool {
	switch c {
	case commandDisable:
	case commandEnable:
	case commandInstall:
	case commandRemove:
	case commandRescan:
	case commandRestart:
	case commandUpdate:
	case commandUpdateNI:
		return true
	}

	return false
}

// CanBeRemote returns true if the associated command can target a remote
// computer.
func (c command) CanBeRemote() bool {
	switch c {
	case commandClasses:
	case commandFind:
	case commandFindAll:
	case commandHwIDs:
	case commandListClass:
	case commandRescan:
	case commandResources:
	case commandSetHwID:
	case commandStack:
	case commandStatus:
		return true
	}

	return false
}
