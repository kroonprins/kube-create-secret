package core

import (
	"fmt"

	"github.com/kroonprins/kube-create-secret/pkg/types"
)

func preResolve(secretTemplate *types.SecretTemplate) (*types.PreResolvedSecretTemplateSpec, error) {
	preResolvedSecretTemplateSpec := types.PreResolvedSecretTemplateSpec{}

	if secretTemplate.Spec.Data != nil {
		preResolvedData, err := PreResolveData(secretTemplate.Spec.Data)
		if err != nil {
			return nil, err
		}
		preResolvedSecretTemplateSpec.PreResolvedData = preResolvedData
	}

	if secretTemplate.Spec.StringData != nil {
		preResolvedStringData, err := PreResolveData(secretTemplate.Spec.StringData)
		if err != nil {
			return nil, err
		}
		preResolvedSecretTemplateSpec.PreResolvedStringData = preResolvedStringData
	}

	if secretTemplate.Spec.Tls != nil {
		preResolvedTls, err := PreResolveTls(secretTemplate.Spec.Tls)
		if err != nil {
			return nil, err
		}
		preResolvedSecretTemplateSpec.PreResolvedTls = &preResolvedTls
	}

	return &preResolvedSecretTemplateSpec, nil
}

func PreResolveData(data interface{}) (interface{}, error) {
	switch typed_data := data.(type) {
	case string:
		return typed_data, nil
	case map[string]interface{}:
		res := make(map[string]interface{})
		for k, v := range typed_data {
			typed_v, ok := v.(string)
			if !ok {
				return nil, fmt.Errorf("unexpected format of value %#v (expected to be string) for %s", v, k)
			}
			res[k] = typed_v
		}
		return res, nil
	default:
		return nil, fmt.Errorf("unexpected format of data: %#v", data)
	}
}

func PreResolveTls(tls *types.Tls) (types.PreResolvedTls, error) {
	toResolve := make(map[string]interface{})
	toResolve["pkcs12"] = tls.Pkcs12
	toResolve["password"] = tls.Password

	return types.PreResolvedTls{
		ToResolve: toResolve,
		Tls:       tls,
	}, nil
}
