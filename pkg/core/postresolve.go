package core

import (
	"fmt"

	"github.com/kroonprins/kube-create-secret/pkg/constants"
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

func PostResolveTls(resolvedTls *types.ResolvedTls) (*types.PostResolvedTls, error) {
	fileContent := resolvedTls.Resolved["pkcs12"].(string)
	password := resolvedTls.Resolved["password"].(string)

	var chainDelimiter string
	if resolvedTls.Tls.CrtConfig != nil {
		chainDelimiter = resolvedTls.Tls.CrtConfig.ChainDelimiter
	}

	key, crt, err := util.ToPEM(fileContent, password, chainDelimiter)
	if err != nil {
		return nil, err
	}
	return &types.PostResolvedTls{
		Key: &types.TlsSecretData{
			Value: key,
			Name:  getKeyName(resolvedTls),
		},
		Crt: &types.TlsSecretData{
			Value: crt,
			Name:  getCrtName(resolvedTls),
		},
	}, nil
}

func getKeyName(resolvedTls *types.ResolvedTls) string {
	if resolvedTls.Tls.KeyConfig != nil && resolvedTls.Tls.KeyConfig.Name != "" {
		return resolvedTls.Tls.KeyConfig.Name
	}

	if resolvedTls.Tls.Name != "" {
		return fmt.Sprintf("%s.key", resolvedTls.Tls.Name)
	}

	return constants.TLS_SECRET_KEY_FIELD
}

func getCrtName(resolvedTls *types.ResolvedTls) string {
	if resolvedTls.Tls.CrtConfig != nil && resolvedTls.Tls.CrtConfig.Name != "" {
		return resolvedTls.Tls.CrtConfig.Name
	}

	if resolvedTls.Tls.Name != "" {
		return fmt.Sprintf("%s.crt", resolvedTls.Tls.Name)
	}

	return constants.TLS_SECRET_CRT_FIELD
}
