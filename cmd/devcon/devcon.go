package main

import (
	"fmt"
	"io/fs"
	"io/ioutil"

	"github.com/mikerourke/go-devcon"
)

func main() {
	dc := devcon.New("")

	output, err := dc.Resources()
	if err != nil {
		fmt.Println(err)
		return
	}

	bytes, err := devcon.MarshalJSON(output)
	if err != nil {
		fmt.Println(err)
		return
	}

	ioutil.WriteFile("test.json", bytes, fs.ModePerm)
}
