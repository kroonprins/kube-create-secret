package read

import (
	"os"

	"github.com/kroonprins/kube-create-secret/pkg/core"
)

type FileReader struct {
}

func NewFileReader() *FileReader {
	return &FileReader{}
}

func (*FileReader) Read(inputFile string, _ core.Config) (bool, []byte, error) {
	if inputFile == "-" {
		return true, nil, nil
	}
	bytes, err := os.ReadFile(inputFile)
	return false, bytes, err
}
