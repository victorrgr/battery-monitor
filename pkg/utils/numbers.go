package utils

import (
	"fmt"
	"strconv"
)

func ParseInt32(str string) (int32, error) {
	parsed, err := strconv.ParseInt(str, 10, 32)
	if err != nil {
		return 0, fmt.Errorf("Error to parse value \"%s\": %v\n", str, err)
	}
	return int32(parsed), nil
}
