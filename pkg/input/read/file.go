package read

import "os"

type FileReader struct {
}

func NewFileReader() *FileReader {
	return &FileReader{}
}

func (*FileReader) Read(inputFile string) (bool, []byte, error) {
	if inputFile == "-" {
		return true, nil, nil
	}
	bytes, err := os.ReadFile(inputFile)
	return false, bytes, err
}
