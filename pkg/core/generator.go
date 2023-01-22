package core

import (
	"encoding/json"
	"fmt"
	"log"

	"github.com/kroonprins/kube-create-secret/pkg/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/klog/v2"
)

func Create(config Config) error {
	secretTemplates, inputFormat, err := read(config)
	if err != nil {
		return err
	}

	config.InputFormat = inputFormat

	secrets := []corev1.Secret{}

	for _, secretTemplate := range secretTemplates {
		err := validate(&secretTemplate)
		if err != nil {
			return err
		}
		if secretTemplate.Spec.APIVersion == "" {
			secretTemplate.Spec.APIVersion = "v1"
		}
		if secretTemplate.Spec.Kind == "" {
			secretTemplate.Spec.Kind = "Secret"
		}
		debug(secretTemplate, "Secret template")

		preResolvedSecretTemplateSpec, err := preResolve(&secretTemplate)
		if err != nil {
			return fmt.Errorf("failed pre-resolution: %v", err)
		}
		debug(preResolvedSecretTemplateSpec, "Pre-resolved")

		resolvedSecretTemplateSpec, err := resolve(preResolvedSecretTemplateSpec)
		if err != nil {
			return fmt.Errorf("failed to resolve data: %v", err)
		}
		debug(resolvedSecretTemplateSpec, "Resolved")

		postResolvedSecretTemplateSpec, err := postResolve(resolvedSecretTemplateSpec)
		if err != nil {
			return fmt.Errorf("failed post-resolution: %v", err)
		}
		debug(postResolvedSecretTemplateSpec, "Post-resolved")

		secret, err := createSecret(&secretTemplate, postResolvedSecretTemplateSpec)
		if err != nil {
			return fmt.Errorf("failed creating secret: %v", err)
		}
		debug(secret, "Secret")

		secrets = append(secrets, *secret)
	}

	return write(config, secrets)
}

func ShowTemplate(config Config) error {
	secretTemplates, inputFormat, err := read(config)
	if err != nil {
		return err
	}
	debug(secretTemplates, "Secret templates")

	config.InputFormat = inputFormat

	return write(config, secretTemplates)
}

func StarterTemplate(config Config, templateType types.StarterTemplateType) error {
	secretTemplate, err := NewStarterTemplate(config, templateType)
	if err != nil {
		return err
	}
	debug(secretTemplate, "Secret template")

	if len(config.OutputFormats) == 0 {
		config.OutputFormats = []types.Format{types.YAML}
	}

	return write(config, []types.SecretTemplate{*secretTemplate})
}

func debug(object interface{}, objectType string) {
	if klog.V(1).Enabled() {
		marshalled, err := json.Marshal(object)
		if err != nil {
			log.Fatalf("Marshal error %v", err)
		}
		klog.Infof("%s: %s", objectType, marshalled)
	}
}
