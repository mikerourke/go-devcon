package devcon

type Device struct {
	Name        string
	Description string
}

func (dc *DevCon) Find() ([]Device, error) {
	lines, err := dc.run(commandFind)
	if err != nil {
		return nil, err
	}

	devices := make([]Device, 0)

	for _, valuePair := range parseColonSeparatedLines(lines) {
		device := Device{
			Name:        valuePair[0],
			Description: valuePair[1],
		}

		devices = append(devices, device)
	}

	return devices, nil
}

func (dc *DevCon) FindAll() ([]Device, error) {
	lines, err := dc.run(commandFindAll)
	if err != nil {
		return nil, err
	}

	devices := make([]Device, 0)

	for _, valuePair := range parseColonSeparatedLines(lines) {
		device := Device{
			Name:        valuePair[0],
			Description: valuePair[1],
		}

		devices = append(devices, device)
	}

	return devices, nil
}
