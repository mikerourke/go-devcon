package devcon_test

import (
	"fmt"

	"github.com/mikerourke/go-devcon"
)

func ExampleDevCon_DriverFiles() {
	dc := devcon.New(`\path\to\devcon.exe`)

	dfs, _ := dc.DriverFiles(`ACPI\PNP0303\4&1D401FB5&0`)

	df := dfs[3]

	fmt.Printf("ID: %s, INF file: %s, INF section: %s\n",
		df.Device.ID, df.INFFile, df.INFSection)
	fmt.Printf("Files: %q", df.Files)
	// Output:
	// ID: ACPI\PNP0303\4&1D401FB5&0, INF file: c:\windows\inf\keyboard.inf, INF section: STANDARD_Inst
	// Files: ["C:\\WINDOWS\\system32\\DRIVERS\\i8042prt.sys" "C:\\WINDOWS\\system32\\DRIVERS\\kbdclass.sys"]
}

func ExampleDevCon_DriverNodes() {
	dc := devcon.New(`\path\to\devcon.exe`)

	dns, _ := dc.DriverNodes("*")

	fmt.Printf("ID: %s, Node description: %s\n",
		dns[0].Device.ID, dns[0].Nodes[0].Description)
	// Output: ID: ACPI\ACPI0010\2&DABA3FF&0, Node description: Generic Bus
}
