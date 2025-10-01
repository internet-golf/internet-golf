package utils

import (
	"encoding/json"
	"fmt"
)

// this is apparently how you have to do this
type JsonObj map[string]any

// utility function to turn a value into json without possibly returning an
// error. should only really be used if it seems incredibly unlikely that
// json.Marshal will panic when given v.
func JsonOrPanic(v any) []byte {
	result, err := json.Marshal(v)
	if err != nil {
		panic(fmt.Sprintf("Could not JSON-serialize value: %v", v))
	}
	return result
}
