package utils

import (
	"encoding/json"
	"fmt"
)

func PrintRecord(label string, record interface{}) {
	data, err := json.MarshalIndent(record, "", "  ")
	if err != nil {
		return
	}
	fmt.Printf("%s\n %s:\n", label, string(data))
}
