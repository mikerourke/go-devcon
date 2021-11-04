package devcon

// Reboot stops and then starts the operating system.
//
// If the user has open files on the computer or a program will not close, the
// system does not reboot until the user has responded to system prompts to
// close the files or end the process.
//
// Cannot be run with the WithRemoteComputer() option.
//
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-reboot for more information.
func (dc *DevCon) Reboot() error {
	_, err := dc.run(commandReboot)

	return err
}
