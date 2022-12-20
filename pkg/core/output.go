package core

import (
	"fmt"
)

var Marshallers []Marshaller
var OutputWriters []OutputWriter

type OutputWriter interface {
	Write([]byte) (bool, error)
}

type Marshaller interface {
	Marshal(config Config, secrets []interface{}) (bool, []byte, error)
}

func write[T interface{}](config Config, items []T) error {
	for _, marshaller := range Marshallers {
		toMarshal := []interface{}{}
		for _, item := range items {
			toMarshal = append(toMarshal, item)
		}
		skipped, bytes, err := marshaller.Marshal(config, toMarshal)
		if err != nil {
			return fmt.Errorf("failed to marshal: %v", err)
		}
		if skipped {
			continue
		}
		for _, writer := range OutputWriters {
			skipped, err := writer.Write(bytes)
			if err != nil {
				return fmt.Errorf("failed to write: %v", err)
			}
			if skipped {
				continue
			}
			break
		}
	}
	return nil
}
