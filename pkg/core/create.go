package core

import (
	"encoding/json"
	"fmt"

	"github.com/kroonprins/kube-create-secret/pkg/constants"
	"github.com/kroonprins/kube-create-secret/pkg/types"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createSecret(secretTemplate *types.SecretTemplate, postResolvedSecretTemplateSpec *types.PostResolvedSecretTemplateSpec, config Config) (*corev1.Secret, error) {
	var data = postResolvedSecretTemplateSpec.PostResolvedData
	stringData := postResolvedSecretTemplateSpec.PostResolvedStringData

	if postResolvedSecretTemplateSpec.PostResolvedTls.Key != nil || postResolvedSecretTemplateSpec.PostResolvedTls.Crt != nil {
		if data == nil {
			data = make(map[string][]byte)
		}
		data[postResolvedSecretTemplateSpec.PostResolvedTls.Key.Name] = postResolvedSecretTemplateSpec.PostResolvedTls.Key.Value
		data[postResolvedSecretTemplateSpec.PostResolvedTls.Crt.Name] = postResolvedSecretTemplateSpec.PostResolvedTls.Crt.Value
	}

	metaData, err := getMetaData(secretTemplate, config)
	if err != nil {
		return nil, fmt.Errorf("failed to construct meta data for result: %v", err)
	}

	return &corev1.Secret{
		TypeMeta:   secretTemplate.Spec.TypeMeta,
		ObjectMeta: *metaData,
		Immutable:  secretTemplate.Spec.Immutable,
		Data:       data,
		StringData: stringData,
		Type:       secretTemplate.Spec.Type,
	}, nil
}

func getMetaData(secretTemplate *types.SecretTemplate, config Config) (*v1.ObjectMeta, error) {
	var (
		bytes []byte
		err   error
	)
	if (len(config.OutputFormats) == 0 && config.InputFormat == types.YAML) ||
		(len(config.OutputFormats) == 1 && config.OutputFormats[0] == types.YAML) {
		// make use of multi-line for yaml output
		bytes, err = json.MarshalIndent(*secretTemplate, "", "  ")
	} else {
		bytes, err = json.Marshal(*secretTemplate)
	}
	if err != nil {
		return nil, err
	}

	res := &v1.ObjectMeta{}
	secretTemplate.Spec.ObjectMeta.DeepCopyInto(res)

	var annotations = make(map[string]string)
	for k, v := range secretTemplate.Spec.ObjectMeta.GetAnnotations() {
		annotations[k] = v
	}
	annotations[constants.ANNOTATION] = string(bytes)
	res.SetAnnotations(annotations)

	return res, nil
}
