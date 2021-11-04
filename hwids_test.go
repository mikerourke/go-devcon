package devcon_test

import (
	"fmt"

	"github.com/mikerourke/go-devcon"
)

func ExampleDevCon_HwIDs_local() {
	dc := devcon.New(`\path\to\devcon.exe`)

	hwids, _ := dc.HwIDs("*")

	for _, hwid := range hwids {
		fmt.Printf("Device ID: %s, Hardware IDs: %s\n", hwid.Device.ID, hwid.HardwareIDs)
	}
	// Output:
	// Device ID: ACPI\ACPI0010\2&DABA3FF&0, Hardware IDs: [ACPI\ACPI0010 *ACPI0010]
	// Device ID: ACPI\FIXEDBUTTON\2&DABA3FF&0, Hardware IDs: [ACPI\FixedButton *FixedButton]
	// ...
}

func ExampleDevCon_HwIDs_remote() {
	dc := devcon.New(`\path\to\devcon.exe`)

	hwids, _ := dc.WithRemoteComputer(`\\server01`).HwIDs("*floppy*")

	for _, hwid := range hwids {
		fmt.Printf("Device ID: %s, Hardware IDs: %s\n", hwid.Device.ID, hwid.HardwareIDs)
	}
	// Output:
	// Device ID: ACPI\ACPI0010\2&DABA3FF&0, Hardware IDs: [ACPI\ACPI0010 *ACPI0010]
	// Device ID: ACPI\FIXEDBUTTON\2&DABA3FF&0, Hardware IDs: [ACPI\FixedButton *FixedButton]
	// ...
}
