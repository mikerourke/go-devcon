package devcon_test

import (
	"fmt"

	"github.com/mikerourke/go-devcon"
)

func ExampleDevCon_Disable_specific() {
	dc := devcon.New(`\path\to\devcon.exe`)

	_ = dc.WithConditionalReboot().Disable(`@USB\ROOT_HUB\4&2A40B465&0`)
}

func ExampleDevCon_Disable_pattern() {
	dc := devcon.New(`\path\to\devcon.exe`)

	_ = dc.Disable("USB*")
}

func ExampleDevCon_Enable_specific() {
	dc := devcon.New(`\path\to\devcon.exe`)

	_ = dc.Enable(`'*PNP0000`)
}

func ExampleDevCon_Enable_class() {
	dc := devcon.New(`\path\to\devcon.exe`)

	_ = dc.WithConditionalReboot().Enable("=Printer")
}

func ExampleDevCon_Find_everything() {
	dc := devcon.New(`\path\to\devcon.exe`)

	devs, _ := dc.Find("*")
	for _, dev := range devs {
		fmt.Printf("ID: %s, Name: %s\n", dev.ID, dev.Name)
	}
	// Output:
	// ID: ACPI\ACPI0010\2&DABA3FF&0, Name: Generic Bus
	// ID: ACPI\FIXEDBUTTON\2&DABA3FF&0, Name: ACPI Fixed Feature Button
	// ...
}

func ExampleDevCon_Find_pattern() {
	dc := devcon.New(`\path\to\devcon.exe`)

	devs, _ := dc.Find(`@root\legacy*`)
	for _, dev := range devs {
		fmt.Printf("ID: %s, Name: %s\n", dev.ID, dev.Name)
	}
	// Output:
	// ID: ROOT\LEGACY_AFD\0000, Name: AFD Networking Support Environment
	// ID: ROOT\LEGACY_BEEP\0000, Name: Beep
	// ...
}

func ExampleDevCon_FindAll() {
	dc := devcon.New(`\path\to\devcon.exe`)

	devs, _ := dc.FindAll("=net")
	for _, dev := range devs {
		fmt.Printf("ID: %s, Name: %s\n", dev.ID, dev.Name)
	}
	// Output:
	// ID: ROOT\MS_L2TPMINIPORT\0000, Name: WAN Miniport (L2TP)
	// ID: ROOT\MS_NDISWANIP\0000, Name: WAN Miniport (IP)
	// ...
}

func ExampleDevCon_Install() {
	dc := devcon.New(`\path\to\devcon.exe`)

	_ = dc.Install(`c:\windows\inf\keyboard.inf`, "*PNP030b")
}

func ExampleDevCon_Remove() {
	dc := devcon.New(`\path\to\devcon.exe`)

	_ = dc.WithConditionalReboot().Remove(`@usb\*`)
}

func ExampleDevCon_Restart() {
	dc := devcon.New(`\path\to\devcon.exe`)

	statuses, _ := dc.Restart("=net", `@'ROOT\*MSLOOP\0000`)

	fmt.Printf("ID: %s, Restarted: %v\n", statuses[0].ID, statuses[0].WasRestarted)
	// Output: ID: ROOT\*MSLOOP\0000, Restarted: true
}
