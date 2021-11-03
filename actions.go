package devcon

func (dc *DevCon) Disable(ids ...string) error {
	_, err := dc.run(commandDisable, ids...)

	return err
}

func (dc *DevCon) Enable(ids ...string) error {
	_, err := dc.run(commandEnable, ids...)

	return err
}

func (dc *DevCon) Update(infFile string, ids ...string) error {
	args := []string{infFile}
	args = append(args, ids...)

	_, err := dc.run(commandUpdate, args...)

	return err
}
