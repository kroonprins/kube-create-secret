package core

import (
	"encoding/json"
	"fmt"

	"github.com/kroonprins/kube-create-secret/pkg/constants"
	"github.com/kroonprins/kube-create-secret/pkg/types"
	corev1 "k8s.io/api/core/v1"
	v1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func createSecret(secretTemplate *types.SecretTemplate, postResolvedSecretTemplateSpec *types.PostResolvedSecretTemplateSpec) (*corev1.Secret, error) {
	var data = postResolvedSecretTemplateSpec.PostResolvedData
	stringData := postResolvedSecretTemplateSpec.PostResolvedStringData

	if postResolvedSecretTemplateSpec.PostResolvedTls.Key != nil || postResolvedSecretTemplateSpec.PostResolvedTls.Crt != nil {
		if data == nil {
			data = make(map[string][]byte)
		}
		data[constants.TLS_SECRET_KEY_FIELD] = postResolvedSecretTemplateSpec.PostResolvedTls.Key
		data[constants.TLS_SECRET_CRT_FIELD] = postResolvedSecretTemplateSpec.PostResolvedTls.Crt
	}

	metaData, err := getMetaData(secretTemplate)
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

func getMetaData(secretTemplate *types.SecretTemplate) (*v1.ObjectMeta, error) {
	bytes, err := json.Marshal(*secretTemplate)
	if err != nil {
		return nil, err
	}

	res := &v1.ObjectMeta{}
	secretTemplate.ObjectMeta.DeepCopyInto(res)

	var annotations = make(map[string]string)
	for k, v := range secretTemplate.ObjectMeta.GetAnnotations() {
		annotations[k] = v
	}
	annotations[constants.ANNOTATION] = string(bytes)
	res.SetAnnotations(annotations)

	return res, nil
}
