package unmarshal

import (
	"encoding/json"
	"fmt"

	"github.com/kroonprins/kube-create-secret/pkg/constants"
	"github.com/kroonprins/kube-create-secret/pkg/types"

	corev1 "k8s.io/api/core/v1"
)

type JsonUnmarshaller struct {
}

func NewJsonUnmarshaller() *JsonUnmarshaller {
	return &JsonUnmarshaller{}
}

func (*JsonUnmarshaller) Unmarshal(bytes []byte) ([]types.SecretTemplate, types.Format, error) {
	secretTemplate := types.SecretTemplate{}
	errSecretTemplate := json.Unmarshal(bytes, &secretTemplate)
	if errSecretTemplate == nil && secretTemplate.Kind == constants.SECRET_TEMPLATE_KIND {
		return []types.SecretTemplate{secretTemplate}, types.JSON, nil
	}

	secretTemplateList := types.SecretTemplateList{}
	errSecretTemplateList := json.Unmarshal(bytes, &secretTemplateList)
	if errSecretTemplateList != nil {
		return nil, 0, fmt.Errorf("failed to unmarshal to either SecretTemplate (%v) or List (%v)", errSecretTemplate, errSecretTemplateList)
	}
	return secretTemplateList.Items, types.JSON, nil
}

type ReCreateJsonUnmarshaller struct {
	jsonUnmarshaller JsonUnmarshaller
}

func NewReCreateJsonUnmarshaller() *ReCreateJsonUnmarshaller {
	return &ReCreateJsonUnmarshaller{
		jsonUnmarshaller: *NewJsonUnmarshaller(),
	}
}

func (reCreateJsonUnmarshaller *ReCreateJsonUnmarshaller) Unmarshal(bytes []byte) ([]types.SecretTemplate, types.Format, error) {
	// try unmarshal to normal secret
	secret := corev1.Secret{}
	errSecret := json.Unmarshal(bytes, &secret)
	if errSecret == nil && secret.Kind == "Secret" {
		if annotation, exists := secret.Annotations[constants.ANNOTATION]; exists {
			unmarshalled, _, err := reCreateJsonUnmarshaller.jsonUnmarshaller.Unmarshal([]byte(annotation))
			return unmarshalled, types.JSON, err
		} else {
			return nil, 0, fmt.Errorf("annotation %s not present", constants.ANNOTATION)
		}
	}

	// try unmarshal to list of normal secrets
	secretList := corev1.SecretList{}
	errSecretList := json.Unmarshal(bytes, &secretList)
	if errSecretList == nil && secret.Kind == "List" && len(secretList.Items) > 0 {
		res := []types.SecretTemplate{}
		for _, secret := range secretList.Items {
			if secret.Kind != "Secret" {
				return nil, 0, fmt.Errorf("list should contain only type Secret")
			}
			if annotation, exists := secret.Annotations[constants.ANNOTATION]; exists {
				unmarshalled, _, err := reCreateJsonUnmarshaller.jsonUnmarshaller.Unmarshal([]byte(annotation))
				if err != nil {
					return nil, 0, fmt.Errorf("unmarshal error for %s: %v", annotation, err)
				}
				res = append(res, unmarshalled...)
			} else {
				return nil, 0, fmt.Errorf("annotation %s not present", constants.ANNOTATION)
			}
		}
		return res, types.JSON, nil
	}

	return nil, 0, fmt.Errorf("failed to unmarshal to either Secret (%v) or List (%v)", errSecret, errSecretList)
}
