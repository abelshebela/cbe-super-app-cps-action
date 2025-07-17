package utils

import "encoding/json"

// TODO: This function needs fix, It's not working as expected but its used in multiple places
func StructToMap(data interface{}) (map[string]interface{}, error) {
	bytes, err := json.Marshal(data)
	if err != nil {
		return nil, err
	}

	var result map[string]interface{}
	err = json.Unmarshal(bytes, &result)
	return result, err
}
