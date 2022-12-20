package unmarshal

import (
	"encoding/json"
	"fmt"

	"github.com/kroonprins/kube-create-secret/pkg/types"

	ssv1 "github.com/bitnami-labs/sealed-secrets/pkg/apis/sealedsecrets/v1alpha1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	yaml "sigs.k8s.io/yaml"
)

type ReCreateJsonSealedSecretUnmarshaller struct {
	reCreateJsonUnmarshaller ReCreateJsonUnmarshaller
}

func NewReCreateJsonSealedSecretUnmarshaller() *ReCreateJsonSealedSecretUnmarshaller {
	return &ReCreateJsonSealedSecretUnmarshaller{
		reCreateJsonUnmarshaller: *NewReCreateJsonUnmarshaller(),
	}
}

func (reCreateJsonSealedSecretUnmarshaller *ReCreateJsonSealedSecretUnmarshaller) Unmarshal(bytes []byte) ([]types.SecretTemplate, types.Format, error) {
	unmarshalled, err := unmarshal(bytes, jsonUnmarshal, reCreateJsonSealedSecretUnmarshaller.reCreateJsonUnmarshaller)
	return unmarshalled, types.JSON, err
}

type ReCreateYamlSealedSecretUnmarshaller struct {
	reCreateJsonUnmarshaller ReCreateJsonUnmarshaller
}

func NewReCreateYamlSealedSecretUnmarshaller() *ReCreateYamlSealedSecretUnmarshaller {
	return &ReCreateYamlSealedSecretUnmarshaller{
		reCreateJsonUnmarshaller: *NewReCreateJsonUnmarshaller(),
	}
}

func (reCreateYamlSealedSecretUnmarshaller *ReCreateYamlSealedSecretUnmarshaller) Unmarshal(bytes []byte) ([]types.SecretTemplate, types.Format, error) {
	res := []types.SecretTemplate{}
	docs, err := splitYAML(bytes)
	if err != nil {
		return nil, 0, err
	}
	for _, doc := range docs {
		secretTemplates, err := unmarshal(doc, yamlUnmarshal, reCreateYamlSealedSecretUnmarshaller.reCreateJsonUnmarshaller)
		if err != nil {
			return nil, 0, err
		}
		res = append(res, secretTemplates...)
	}
	return res, types.YAML, nil
}

var yamlUnmarshal = func(bytes []byte, o interface{}) error {
	return yaml.Unmarshal(bytes, o)
}

var jsonUnmarshal = func(bytes []byte, o interface{}) error {
	return json.Unmarshal(bytes, o)
}

func unmarshal(bytes []byte, unmarshaller func([]byte, interface{}) error, recreateUnmarshaller ReCreateJsonUnmarshaller) ([]types.SecretTemplate, error) {
	// try unmarshal to normal sealed secret
	sealedSecret := ssv1.SealedSecret{}
	errSealedSecret := unmarshaller(bytes, &sealedSecret)
	if errSealedSecret == nil && sealedSecret.Kind == "SealedSecret" {
		unmarshalled, err := unmarshalSealedSecret(sealedSecret, recreateUnmarshaller)
		return unmarshalled, err
	}

	// try unmarshal to list of normal sealed secrets
	sealedSecretList := ssv1.SealedSecretList{}
	errSealedSecretList := yaml.Unmarshal(bytes, &sealedSecretList)
	if errSealedSecretList == nil && sealedSecretList.Kind == "List" && len(sealedSecretList.Items) > 0 {
		res := []types.SecretTemplate{}
		for _, sealedSecret := range sealedSecretList.Items {
			if sealedSecret.Kind != "SealedSecret" {
				return nil, fmt.Errorf("list should contain only type SealedSecret")
			}
			unmarshalled, err := unmarshalSealedSecret(sealedSecret, recreateUnmarshaller)
			if err != nil {
				return nil, err
			}
			res = append(res, unmarshalled...)

		}
		return res, nil
	}

	return nil, fmt.Errorf("failed to unmarshal to either SealedSecret (%v) or List (%v)", errSealedSecret, errSealedSecretList)
}

func unmarshalSealedSecret(sealedSecret ssv1.SealedSecret, unmarshaller ReCreateJsonUnmarshaller) ([]types.SecretTemplate, error) {
	secret := corev1.Secret{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "v1",
			Kind:       "Secret",
		},
		ObjectMeta: sealedSecret.Spec.Template.ObjectMeta,
	}
	bytes, err := json.Marshal(secret)
	if err != nil {
		return nil, err
	}
	unmarshalled, _, err := unmarshaller.Unmarshal(bytes)
	if err != nil {
		return nil, fmt.Errorf("unmarshal error for %s: %v", bytes, err)
	}
	return unmarshalled, err
}
