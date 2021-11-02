package devcon

import "regexp"

var (
	reName            = regexp.MustCompile(`Name: (?P<Name>.*)`)
	reNameDesc        = regexp.MustCompile(`(?P<Name>.*): (?P<Desc>.*)`)
	reDriverInstalled = regexp.MustCompile(`Driver installed from (?P<INFFile>.*) \[(?P<INFSection>.*)].*`)
	reDriverNoInfo    = regexp.MustCompile(`No driver information`)
	reDriverFilePath  = regexp.MustCompile(`C:\\(.*)`)
	reDriverNode      = regexp.MustCompile(`DriverNode #(.*):`)
	reFieldIsValue    = regexp.MustCompile(`(?P<Field>.*) is (?P<Value>.*)`)
	reFieldAreValue   = regexp.MustCompile(`(?P<Field>.*) are (?P<Value>.*)`)
	reHash            = regexp.MustCompile(`#`)
)
