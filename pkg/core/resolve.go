package core

import (
	"fmt"

	"github.com/kroonprins/kube-create-secret/pkg/types"
	"github.com/kroonprins/vals"
)

func resolve(preResolvedSecretTemplateSpec *types.PreResolvedSecretTemplateSpec) (*types.ResolvedSecretTemplateSpec, error) {
	resolvedSecretTemplateSpec := types.ResolvedSecretTemplateSpec{}

	// TODO get options via cmd line args
	runtime, err := vals.New(vals.Options{})
	if err != nil {
		return nil, fmt.Errorf("error initializing vals: %v", err)
	}

	if preResolvedSecretTemplateSpec.PreResolvedData != nil {
		resolved, err := runtime.Eval(map[string]interface{}{"inline": preResolvedSecretTemplateSpec.PreResolvedData})
		if err != nil {
			return nil, fmt.Errorf("failure resolving vals for data: %v", err)
		}
		resolvedSecretTemplateSpec.ResolvedData = resolved["inline"].(map[string]interface{})
	}

	if preResolvedSecretTemplateSpec.PreResolvedStringData != nil {
		resolved, err := runtime.Eval(map[string]interface{}{"inline": preResolvedSecretTemplateSpec.PreResolvedStringData})
		if err != nil {
			return nil, fmt.Errorf("failure resolving vals for string data: %v", err)
		}
		resolvedSecretTemplateSpec.ResolvedStringData = resolved["inline"].(map[string]interface{})
	}

	if preResolvedSecretTemplateSpec.PreResolvedTls != nil {
		resolved, err := runtime.Eval(map[string]interface{}{"inline": preResolvedSecretTemplateSpec.PreResolvedTls.ToResolve})
		if err != nil {
			return nil, fmt.Errorf("failure resolving vals for tls: %v", err)
		}
		resolvedSecretTemplateSpec.ResolvedTls = &types.ResolvedTls{
			Resolved: resolved["inline"].(map[string]interface{}),
			Tls:      preResolvedSecretTemplateSpec.PreResolvedTls.Tls,
		}
	}

	return &resolvedSecretTemplateSpec, nil
}
