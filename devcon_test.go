package devcon_test

import (
	"fmt"
	"strings"

	"github.com/mikerourke/go-devcon"
)

func Example_install() {
	dc := devcon.New(`C:\windows\devcon.exe`)

	devs, err := dc.Find("*VEN_8086*")
	if err != nil {
		fmt.Println(err)

		return
	}

	var matchingDev devcon.Device

	for _, dev := range devs {
		if strings.Contains(dev.Name, "PRO") {
			matchingDev = dev
			break
		}
	}

	if matchingDev.ID != "" {
		err = dc.Install(`C:\drivers\PRO1000\WIN32\E1000325.INF`, matchingDev.ID)
		if err != nil {
			fmt.Printf("error installing: %s", err)
		} else {
			fmt.Println("Successfully installed")

		}
	}
}
