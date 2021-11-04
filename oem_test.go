package devcon_test

import (
	"fmt"

	"github.com/mikerourke/go-devcon"
)

func ExampleDevCon_DPEnum() {
	dc := devcon.New(`\path\to\devcon.exe`)

	pkgs, _ := dc.DPEnum()

	fmt.Printf("Name: %s, Provider: %s, Class: %s",
		pkgs[0].Name, pkgs[0].Provider, pkgs[0].Class)
	// Output: Name: oem2.inf, Provider: Microsoft, Class: unknown
}
