package core

import (
	"fmt"

	"github.com/kroonprins/kube-create-secret/pkg/constants"
	"github.com/kroonprins/kube-create-secret/pkg/types"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func NewStarterTemplate(config Config, templateType types.StarterTemplateType) (*types.SecretTemplate, error) {
	if templateType == types.DATA {
		return &types.SecretTemplate{
			TypeMeta: metav1.TypeMeta{
				Kind:       constants.SECRET_TEMPLATE_KIND,
				APIVersion: constants.SECRET_TEMPLATE_API_VERSION,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "[insert template name]",
			},
			Spec: types.SecretTemplateSpec{
				Secret: corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "[insert secret name]",
						Namespace: "[insert namespace] (optional)",
					},
					Type: corev1.SecretTypeOpaque,
				},
				Data: map[string]string{
					"[insert key]": "ref+[insert provider]://[insert provider config]",
				},
			},
		}, nil
	} else if templateType == types.STRINGDATA {
		return &types.SecretTemplate{
			TypeMeta: metav1.TypeMeta{
				Kind:       constants.SECRET_TEMPLATE_KIND,
				APIVersion: constants.SECRET_TEMPLATE_API_VERSION,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "[insert template name]",
			},
			Spec: types.SecretTemplateSpec{
				Secret: corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "[insert secret name]",
						Namespace: "[insert namespace] (optional)",
					},
					Type: corev1.SecretTypeOpaque,
				},
				StringData: map[string]string{
					"[insert key]": "ref+[insert provider]://[insert provider config]",
				},
			},
		}, nil
	} else if templateType == types.TLS {
		return &types.SecretTemplate{
			TypeMeta: metav1.TypeMeta{
				Kind:       constants.SECRET_TEMPLATE_KIND,
				APIVersion: constants.SECRET_TEMPLATE_API_VERSION,
			},
			ObjectMeta: metav1.ObjectMeta{
				Name: "[insert template name]",
			},
			Spec: types.SecretTemplateSpec{
				Secret: corev1.Secret{
					ObjectMeta: metav1.ObjectMeta{
						Name:      "[insert secret name]",
						Namespace: "[insert namespace] (optional)",
					},
					Type: corev1.SecretTypeTLS,
				},
				Tls: &types.Tls{
					Pkcs12:   "ref+[insert provider]://[insert provider config]",
					Password: "ref+[insert provider]://[insert provider config] (optional)",
					Name:     "[insert name] (optional)",
					KeyConfig: &types.TlsKeyConfig{
						Name: "[insert name] (optional)",
					},
					CrtConfig: &types.TlsCrtConfig{
						Name:           "[insert name (optional)]",
						ChainDelimiter: "[insert delimiter (optional)]",
					},
				},
			},
		}, nil
	}
	return nil, fmt.Errorf("unhandled template type %s", types.StarterTemplateTypes[templateType][0])
}
