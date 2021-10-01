package serialization

import (
	"encoding/json"
	"fmt"
)

// Converts a struct into a json string. If input is a string and not a struct
// no conversion will take place and string is returned instead
func ToJSON(in interface{}) (string, error) {
	// check if input is pure string
	if _, ok := in.(string); ok {
		return fmt.Sprintf("%v", in), nil // do no conversion return pure string
	} else {
		byteArray, err := json.Marshal(in)
		return string(byteArray), err
	}
}
