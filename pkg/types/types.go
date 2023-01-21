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
	// TODO: support alias of entry from keystore to select?
	Pkcs12    string        `json:"pkcs12,omitempty"`
	Password  string        `json:"password,omitempty"`
	Name      string        `json:"name,omitempty"`
	KeyConfig *TlsKeyConfig `json:"key,omitempty"`
	CrtConfig *TlsCrtConfig `json:"crt,omitempty"`
}

type TlsKeyConfig struct {
	Name string `json:"name,omitempty"`
}

type TlsCrtConfig struct {
	Name           string `json:"name,omitempty"`
	ChainDelimiter string `json:"delimiter,omitempty"`
}

type PreResolvedTls struct {
	ToResolve map[string]interface{}
	Tls       *Tls
}

type ResolvedTls struct {
	Resolved map[string]interface{}
	Tls      *Tls
}

type PostResolvedTls struct {
	Key *TlsSecretData
	Crt *TlsSecretData
}

type TlsSecretData struct {
	Value []byte
	Name  string
}

type PreResolvedSecretTemplateSpec struct {
	PreResolvedData       interface{}
	PreResolvedStringData interface{}
	PreResolvedTls        *PreResolvedTls
}

type ResolvedSecretTemplateSpec struct {
	ResolvedData       map[string]interface{}
	ResolvedStringData map[string]interface{}
	ResolvedTls        *ResolvedTls
}

type PostResolvedSecretTemplateSpec struct {
	PostResolvedData       map[string][]byte
	PostResolvedStringData map[string]string
	PostResolvedTls        PostResolvedTls
}
