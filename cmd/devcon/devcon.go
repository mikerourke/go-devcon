package main

import (
	"bytes"
	"encoding/json"
	"fmt"

	"github.com/mikerourke/go-devcon"
)

func main() {
	dc := devcon.New("")

	results, err := dc.Restart("*")
	if err != nil {
		fmt.Println(err)

		return
	}

	fmt.Printf("ID: %s, Was restarted: %v\n", results[0].ID, results[0].WasRestarted)
}

// marshalUnescapedJSON returns the JSON representation of the specified interface
// without HTML escaped.
//nolint:deadcode,unused // This is for testing purposes.
func marshalUnescapedJSON(t interface{}) ([]byte, error) {
	buffer := &bytes.Buffer{}

	encoder := json.NewEncoder(buffer)
	encoder.SetEscapeHTML(false)

	err := encoder.Encode(t)

	return buffer.Bytes(), err
}
