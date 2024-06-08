package utils

import (
	"encoding/json"
	"fmt"
	"reflect"
	"strings"
)

// CompareJSON compares two JSON strings and returns the differences as a string.
func CompareJSON(json1, json2 string) string {
	// Unmarshal JSON strings into maps
	var obj1, obj2 map[string]interface{}

	if err := json.Unmarshal([]byte(json1), &obj1); err != nil {
		return fmt.Sprintf("Error parsing Current Value: %s", err)
	}
	if err := json.Unmarshal([]byte(json2), &obj2); err != nil {
		return fmt.Sprintf("Error parsing Prior Value: %s", err)
	}

	// Compare the maps
	return compareMaps(obj1, obj2, "")
}

// compareMaps recursively compares two maps and returns the differences as a string.
func compareMaps(m1, m2 map[string]interface{}, prefix string) string {
	var diff strings.Builder
	for key, val1 := range m1 {
		if val2, ok := m2[key]; ok {
			switch v1 := val1.(type) {
			case map[string]interface{}:
				if v2, ok := val2.(map[string]interface{}); ok {
					// Recursively compare nested maps
					diff.WriteString(compareMaps(v1, v2, fmt.Sprintf("%s%s.", prefix, key)))
				} else {
					diff.WriteString(fmt.Sprintf("%s%s: type mismatch\n", prefix, key))
				}
			case []interface{}:
				// Ignore array comparison for simplicity in this example
			default:
				// Compare values
				if !reflect.DeepEqual(val1, val2) {
					diff.WriteString(fmt.Sprintf("%s%s: %v   ----->   %v     \n", prefix, key, val2, val1))
				}
			}
		} else {
			diff.WriteString(fmt.Sprintf("%s%s: missing in Prior Value\n", prefix, key))
		}
	}

	// Check for keys in JSON2 that are missing in JSON1
	for key := range m2 {
		if _, ok := m1[key]; !ok {
			diff.WriteString(fmt.Sprintf("%s%s: missing in Current Value\n", prefix, key))
		}
	}

	return diff.String()
}
