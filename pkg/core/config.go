package core

import (
	"io"

	"github.com/kroonprins/kube-create-secret/pkg/types"
)

type Config struct {
	InputFiles    []string
	InputReader   io.Reader // will be used if InputFiles contains "-"
	OutputFormats []types.Format
	InputFormat   types.Format // best effort determination of the input format, if there are multiple inputs with different formats then one of them is selected randomly

	// extra configuration for specific providers/outputs/...
	Extra map[string]interface{}
}

func NewConfig() *Config {
	return &Config{
		Extra: make(map[string]interface{}),
	}
}
