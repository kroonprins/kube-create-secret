package write

import (
	"fmt"
)

type StdOutWriter struct {
}

func NewStdOutWriter() *StdOutWriter {
	return &StdOutWriter{}
}

func (*StdOutWriter) Write(bytes []byte) (bool, error) {
	fmt.Print(string(bytes))
	return false, nil
}
