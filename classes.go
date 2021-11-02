package devcon

type Class struct {
	Name        string
	Description string
}

// Classes lists all device setup classes, including classes that devices on the
// system do not use. Valid on local and remote computers.
// See https://docs.microsoft.com/en-us/windows-hardware/drivers/devtest/devcon-classes
func (dc *DevCon) Classes() ([]Class, error) {
	lines, err := dc.runWithoutArgs(commandClasses)
	if err != nil {
		return nil, err
	}

	classes := make([]Class, 0)

	for _, valuePair := range parseColonSeparatedLines(lines) {
		class := Class{
			Name:        valuePair[0],
			Description: valuePair[1],
		}

		classes = append(classes, class)
	}

	return classes, nil
}
