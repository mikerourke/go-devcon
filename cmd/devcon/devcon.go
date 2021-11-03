package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io/ioutil"

	"github.com/mikerourke/go-devcon"
)

func main() {
	dc := devcon.New("")

	output, err := dc.DPEnum()
	if err != nil {
		fmt.Println(err)
		return
	}

	data, err := marshalUnescapedJSON(output)
	if err != nil {
		fmt.Println(err)
		return
	}

	ioutil.WriteFile("out/test2.json", data, 0777)
}

// marshalUnescapedJSON returns the JSON representation of the specified interface
// without HTML escaped.
func marshalUnescapedJSON(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}

	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(t)

	return buffer.Bytes(), err
}
