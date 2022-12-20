package read

import (
	"io"
	"os"
)

type StdInReader struct {
}

func NewStdInReader() *StdInReader {
	return &StdInReader{}
}

func (*StdInReader) Read(inputFile string) (bool, []byte, error) {
	if inputFile != "-" {
		return true, nil, nil
	}
	bytes, err := io.ReadAll(os.Stdin)
	return false, bytes, err
}
