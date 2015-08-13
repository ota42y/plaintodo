package util

import (
	"fmt"
	"strconv"
)

// GetIntAttribute return Atoi(attributes[name])
// If not key exist or fail Atoi, return error
func GetIntAttribute(name string, attributes map[string]string) (int, error) {
	str, ok := attributes[name]
	if !ok {
		return -1, fmt.Errorf("not exist :%s\n", name)
	}

	num, err := strconv.Atoi(str)
	if err != nil {
		return -1, err
	}

	return num, nil
}
