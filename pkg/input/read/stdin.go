package read

import (
	"io"

	"github.com/kroonprins/kube-create-secret/pkg/core"
)

type StdInReader struct {
}

func NewStdInReader() *StdInReader {
	return &StdInReader{}
}

func (*StdInReader) Read(inputFile string, config core.Config) (bool, []byte, error) {
	if inputFile != "-" {
		return true, nil, nil
	}
	bytes, err := io.ReadAll(config.InputReader)
	return false, bytes, err
}
