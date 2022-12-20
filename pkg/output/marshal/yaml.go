package marshal

import (
	"github.com/kroonprins/kube-create-secret/pkg/core"
	"github.com/kroonprins/kube-create-secret/pkg/types"
	yaml "sigs.k8s.io/yaml"
)

var yamlSeparator = []byte("---\n")

type YamlMarshaller struct {
}

func NewYamlMarshaller() *YamlMarshaller {
	return &YamlMarshaller{}
}

func (*YamlMarshaller) Marshal(config core.Config, items []interface{}) (bool, []byte, error) {
	if !(len(config.OutputFormats) == 0 && config.InputFormat == types.YAML) &&
		!(len(config.OutputFormats) == 1 && config.OutputFormats[0] == types.YAML) {
		return true, nil, nil
	}

	res, err := marshalYaml(items)
	if err != nil {
		return false, nil, err
	}
	return false, res, nil
}

func marshalYaml[T interface{}](items []T) ([]byte, error) {
	res := []byte{}
	for _, item := range items {
		bytes, err := yaml.Marshal(item)
		if err != nil {
			return nil, err
		}
		if len(items) > 1 {
			res = append(res, yamlSeparator...)
		}
		res = append(res, bytes...)
	}
	return res, nil
}
