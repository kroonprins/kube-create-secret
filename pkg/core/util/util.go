package util

import "fmt"

var ToBytes = func(r interface{}) ([]byte, error) {
	switch v := r.(type) {
	case string:
		return []byte(v), nil
	case []byte:
		return v, nil
	default:
		return nil, fmt.Errorf("failed to convert '%v' to bytes", r)
	}
}

var ToString = func(r interface{}) (string, error) {
	switch v := r.(type) {
	case string:
		return v, nil
	case []byte:
		return string(v), nil
	default:
		return "", fmt.Errorf("failed to convert '%v' to bytes", r)
	}
}
