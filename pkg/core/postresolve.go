package core

import (
	"github.com/kroonprins/kube-create-secret/pkg/core/util"
	"github.com/kroonprins/kube-create-secret/pkg/types"
)

func postResolve(resolvedSecretTemplateSpec *types.ResolvedSecretTemplateSpec) (*types.PostResolvedSecretTemplateSpec, error) {
	postResolvedSecretTemplateSpec := types.PostResolvedSecretTemplateSpec{}

	if resolvedSecretTemplateSpec.ResolvedData != nil {
		postResolvedData, err := PostResolveData(resolvedSecretTemplateSpec.ResolvedData, util.ToBytes)
		if err != nil {
			return nil, err
		}
		postResolvedSecretTemplateSpec.PostResolvedData = postResolvedData
	}

	if resolvedSecretTemplateSpec.ResolvedStringData != nil {
		postResolvedStringData, err := PostResolveData(resolvedSecretTemplateSpec.ResolvedStringData, util.ToString)
		if err != nil {
			return nil, err
		}
		postResolvedSecretTemplateSpec.PostResolvedStringData = postResolvedStringData
	}

	if resolvedSecretTemplateSpec.ResolvedTls != nil {
		postResolvedTls, err := PostResolveTls(resolvedSecretTemplateSpec.ResolvedTls)
		if err != nil {
			return nil, err
		}
		postResolvedSecretTemplateSpec.PostResolvedTls = *postResolvedTls
	}

	return &postResolvedSecretTemplateSpec, nil
}

func PostResolveData[T []byte | string](data map[string]interface{}, mapper func(interface{}) (T, error)) (map[string]T, error) {
	res := make(map[string]T)
	for k, v := range data {
		mapped, err := mapper(v)
		if err != nil {
			return nil, err
		}
		res[k] = mapped
	}
	return res, nil
}

func PostResolveTls(resolvedTls map[string]interface{}) (*types.ResolvedTls, error) {
	fileContent := resolvedTls["pkcs12"].(string)
	password := resolvedTls["password"].(string)

	key, crt, err := util.ToPEM(fileContent, password)
	if err != nil {
		return nil, err
	}
	return &types.ResolvedTls{
		Key: key,
		Crt: crt,
	}, nil
}
