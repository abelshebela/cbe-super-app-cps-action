package utils

import (
	"bytes"
	"encoding/json"
	"fmt"
)

func PrintJSON(v any) error {
	var b []byte
	var err error

	switch t := v.(type) {
	case []byte:
		b = t
	case string:
		b = []byte(t)
	default:
		b, err = json.Marshal(t)
		if err != nil {
			return err
		}
	}

	var out bytes.Buffer
	if err := json.Indent(&out, b, "", "  "); err != nil {
		return err
	}

	fmt.Println(out.String())
	return nil
}
