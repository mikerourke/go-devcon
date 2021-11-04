package devcon_test

import (
	"fmt"

	"github.com/mikerourke/go-devcon"
)

func ExampleDevCon_Classes() {
	dc := devcon.New(`\path\to\devcon.exe`)

	classes, _ := dc.Classes()

	fmt.Printf("Name: %s, Description: %s", classes[0].Name, classes[0].Description)
	// Output: Name: WCEUSBS, Description: Windows CE USB Devices
}

func ExampleDevCon_ClassFilter_query() {
	dc := devcon.New(`\path\to\devcon.exe`)

	result, _ := dc.ClassFilter("DiskDrive", devcon.ClassFilterUpper)

	fmt.Printf("Filters: %q, Was changed: %v, Requires reboot: %v",
		result.Filters, result.WasChanged, result.RequiresReboot)
	// Output: Filters: ["PartMgr" "Disklog"], Was changed: false, Requires reboot: false
}

func ExampleDevCon_ClassFilter_update() {
	dc := devcon.New(`\path\to\devcon.exe`)

	result, _ := dc.ClassFilter("DiskDrive", devcon.ClassFilterUpper, "+Disklog")

	fmt.Printf("Filters: %q, Was changed: %v, Requires reboot: %v",
		result.Filters, result.WasChanged, result.RequiresReboot)
	// Output: Filters: ["PartMgr" "Disklog"], Was changed: true, Requires reboot: true
}

func ExampleDevCon_ListClass() {
	dc := devcon.New(`\path\to\devcon.exe`)

	classes, _ := dc.ListClass("ports", "keyboard")

	fmt.Printf("Class: 'ports', First Device ID: %s\n", classes["ports"][0].ID)
	fmt.Printf("Class: 'keyboard', First Device Name: %s\n", classes["ports"][0].Name)
	// Output:
	//
	// Class: 'ports', First Device ID: "ACPI\\PNP0400\\1"
	// Class: 'keyboard', First Device Name: "Standard 101/102-Key or Microsoft Natural PS/2 Keyboard"
}
