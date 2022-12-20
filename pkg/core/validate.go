package core

import (
	"fmt"

	"github.com/kroonprins/kube-create-secret/pkg/constants"
	"github.com/kroonprins/kube-create-secret/pkg/types"
)

func validate(secretTemplate *types.SecretTemplate) error {
	if secretTemplate.APIVersion == "" {
		return fmt.Errorf("missing 'apiVersion', should be %s", constants.SECRET_TEMPLATE_API_VERSION)
	}
	if secretTemplate.APIVersion != constants.SECRET_TEMPLATE_API_VERSION {
		return fmt.Errorf("unexpected 'apiVersion' '%s', should be %s", secretTemplate.APIVersion, constants.SECRET_TEMPLATE_API_VERSION)
	}

	if secretTemplate.Kind == "" {
		return fmt.Errorf("missing 'kind', should be %s", constants.SECRET_TEMPLATE_KIND)
	}
	if secretTemplate.Kind != constants.SECRET_TEMPLATE_KIND {
		return fmt.Errorf("unexpected 'kind' '%s', should be %s", secretTemplate.Kind, constants.SECRET_TEMPLATE_KIND)
	}

	if secretTemplate.ObjectMeta.Name == "" {
		return fmt.Errorf("missing 'metaData.name'")
	}

	if secretTemplate.Spec.ObjectMeta.Name == "" {
		return fmt.Errorf("missing 'spec.metadata.name'")
	}
	return nil
}
