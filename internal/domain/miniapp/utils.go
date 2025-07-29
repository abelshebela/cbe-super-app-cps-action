package miniapp

import (
	"encoding/json"
)

// bindAction marshals the source to JSON and unmarshals it into the target.
// Target should be a pointer to the desired struct.
func bindAction(source any, target any) error {
	bytes, err := json.Marshal(source)
	if err != nil {
		return err
	}
	return json.Unmarshal(bytes, target)
}
