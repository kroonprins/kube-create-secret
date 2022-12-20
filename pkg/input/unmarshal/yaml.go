package unmarshal

import (
	"bytes"
	"fmt"
	"io"

	"github.com/kroonprins/kube-create-secret/pkg/constants"
	"github.com/kroonprins/kube-create-secret/pkg/types"
	goyaml "gopkg.in/yaml.v2"
	corev1 "k8s.io/api/core/v1"
	yaml "sigs.k8s.io/yaml"
)

type YamlUnmarshaller struct {
}

func NewYamlUnmarshaller() *YamlUnmarshaller {
	return &YamlUnmarshaller{}
}

func (*YamlUnmarshaller) Unmarshal(bytes []byte) ([]types.SecretTemplate, types.Format, error) {
	res := []types.SecretTemplate{}
	docs, err := splitYAML(bytes)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to split yaml %v", err)
	}
	for _, doc := range docs {
		secretTemplate := types.SecretTemplate{}

		errSecretTemplate := yaml.Unmarshal(doc, &secretTemplate)
		if errSecretTemplate == nil && secretTemplate.Kind == constants.SECRET_TEMPLATE_KIND {
			res = append(res, secretTemplate)
		}

		secretTemplateList := types.SecretTemplateList{}
		errSecretTemplateList := yaml.Unmarshal(doc, &secretTemplateList)
		if errSecretTemplateList != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal to either SecretTemplate (%v) or List (%v)", errSecretTemplate, errSecretTemplateList)
		}
		res = append(res, secretTemplateList.Items...)
	}
	if len(res) == 0 {
		return nil, 0, fmt.Errorf("no items unmarshalled to yaml")
	}
	return res, types.YAML, nil
}

type ReCreateYamlUnmarshaller struct {
	jsonUnmarshaller JsonUnmarshaller
}

func NewReCreateYamlUnmarshaller() *ReCreateYamlUnmarshaller {
	return &ReCreateYamlUnmarshaller{
		jsonUnmarshaller: *NewJsonUnmarshaller(),
	}
}

func (reCreateYamlUnmarshaller *ReCreateYamlUnmarshaller) Unmarshal(bytes []byte) ([]types.SecretTemplate, types.Format, error) {
	res := []types.SecretTemplate{}
	docs, err := splitYAML(bytes)
	if err != nil {
		return nil, 0, fmt.Errorf("failed to split yaml %v", err)
	}
	for _, doc := range docs {
		// try unmarshal to normal secret
		secret := corev1.Secret{}
		errSecret := yaml.Unmarshal(doc, &secret)
		if errSecret == nil && secret.Kind == "Secret" {
			if annotation, exists := secret.Annotations[constants.ANNOTATION]; exists {
				unmarshalled, _, err := reCreateYamlUnmarshaller.jsonUnmarshaller.Unmarshal([]byte(annotation))
				if err != nil {
					return nil, 0, err
				}
				res = append(res, unmarshalled...)
			} else {
				return nil, 0, fmt.Errorf("annotation %s not present", constants.ANNOTATION)
			}
		}

		// try unmarshal to list of normal secrets
		secretList := corev1.SecretList{}
		errSecretList := yaml.Unmarshal(doc, &secretList)
		if errSecretList == nil && secret.Kind == "List" && len(secretList.Items) > 0 {
			for _, secret := range secretList.Items {
				if secret.Kind != "Secret" {
					return nil, 0, fmt.Errorf("list should contain only type Secret")
				}
				if annotation, exists := secret.Annotations[constants.ANNOTATION]; exists {
					unmarshalled, _, err := reCreateYamlUnmarshaller.jsonUnmarshaller.Unmarshal([]byte(annotation))
					if err != nil {
						return nil, 0, fmt.Errorf("unmarshal error for %s: %v", annotation, err)
					}
					res = append(res, unmarshalled...)
				} else {
					return nil, 0, fmt.Errorf("annotation %s not present", constants.ANNOTATION)
				}

			}
			return res, types.YAML, nil
		}

		if errSecret != nil && errSecretList != nil {
			return nil, 0, fmt.Errorf("failed to unmarshal to either Secret (%v) or List (%v)", errSecret, errSecretList)
		}
	}
	if len(res) == 0 {
		return nil, 0, fmt.Errorf("no items unmarshalled to yaml")
	}
	return res, types.YAML, nil
}

func splitYAML(resources []byte) ([][]byte, error) {
	dec := goyaml.NewDecoder(bytes.NewReader(resources))

	var res [][]byte
	for {
		var value interface{}
		err := dec.Decode(&value)
		if err == io.EOF {
			break
		}
		if err != nil {
			return nil, err
		}
		valueBytes, err := goyaml.Marshal(value)
		if err != nil {
			return nil, err
		}
		res = append(res, valueBytes)
	}
	return res, nil
}
