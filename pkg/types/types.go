package types

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

type Format int

const (
	YAML Format = iota
	JSON
	SEALED_SECRET
)

var FormatIds = map[Format][]string{
	YAML:          {"yaml"},
	JSON:          {"json"},
	SEALED_SECRET: {"sealed-secret"},
}

type SecretTemplate struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`
	Spec              SecretTemplateSpec `json:"spec,omitempty"`
}

type SecretTemplateList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []SecretTemplate `json:"items"`
}

type SecretTemplateSpec struct {
	// Standard Secret fields
	corev1.Secret `json:",inline"`

	// Overriden standard Secret fields
	Data       interface{} `json:"data,omitempty"`
	StringData interface{} `json:"stringData,omitempty"`

	// Template-specific fields
	Tls *Tls `json:"tls,omitempty"`
}

type Tls struct {
	Pkcs12   string `json:"pkcs12,omitempty"`
	Password string `json:"password,omitempty"`
}

type ResolvedTls struct {
	Key []byte
	Crt []byte
}

type PreResolvedSecretTemplateSpec struct {
	PreResolvedData       interface{}
	PreResolvedStringData interface{}
	PreResolvedTls        interface{}
}

type ResolvedSecretTemplateSpec struct {
	ResolvedData       map[string]interface{}
	ResolvedStringData map[string]interface{}
	ResolvedTls        map[string]interface{}
}

type PostResolvedSecretTemplateSpec struct {
	PostResolvedData       map[string][]byte
	PostResolvedStringData map[string]string
	PostResolvedTls        ResolvedTls
}
