package core

import (
	"fmt"

	"github.com/kroonprins/kube-create-secret/pkg/types"
)

var InputReaders []InputReader
var Unmarshallers []Unmarshaller

type InputReader interface {
	Read(inputFile string) (bool, []byte, error)
}

type Unmarshaller interface {
	Unmarshal([]byte) ([]types.SecretTemplate, types.Format, error)
}

func read(config Config) ([]types.SecretTemplate, types.Format, error) {
	res := []types.SecretTemplate{}
	var inputFormat types.Format
	for _, inputFile := range config.InputFiles {
		for _, inputReader := range InputReaders {
			skipped, bytes, err := inputReader.Read(inputFile)
			if err != nil {
				return nil, 0, fmt.Errorf("unable to read file '%s': %v", inputFile, err)
			}
			if skipped {
				continue
			}
			errors := []error{}
			success := false
			for _, unmarshaller := range Unmarshallers {
				unmarshalled, format, err := unmarshaller.Unmarshal(bytes)
				if err != nil {
					errors = append(errors, err)
					continue
				}
				res = append(res, unmarshalled...)
				inputFormat = format
				success = true
				break
			}
			if !success {
				return nil, 0, fmt.Errorf("could not parse input of %s: %v", inputFile, errors)
			}
		}
	}
	return res, inputFormat, nil
}
