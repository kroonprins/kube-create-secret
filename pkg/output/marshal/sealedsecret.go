package marshal

import (
	b "bytes"
	"context"
	"crypto/rsa"
	"fmt"
	"io"
	"log"

	ssv1 "github.com/bitnami-labs/sealed-secrets/pkg/apis/sealedsecrets/v1alpha1"
	"github.com/bitnami-labs/sealed-secrets/pkg/kubeseal"
	"github.com/kroonprins/kube-create-secret/pkg/core"
	"github.com/kroonprins/kube-create-secret/pkg/types"
	corev1 "k8s.io/api/core/v1"
	"k8s.io/client-go/kubernetes/scheme"
	"k8s.io/client-go/tools/clientcmd"
	yaml "sigs.k8s.io/yaml"
)

const SEALED_SECRET_EXTRA_CONFIG_KEY = "kubeseal-config"

type KubesealConfig struct {
	CertURL        string
	ControllerNs   string
	ControllerName string
	Kubeconfig     string
	AllowEmptyData bool
	SealingScope   ssv1.SealingScope

	ConfigOverrides clientcmd.ConfigOverrides
}

type SealedSecretMarshaller struct {
}

func NewSealedSecretMarshaller() *SealedSecretMarshaller {
	return &SealedSecretMarshaller{}
}

func (sealedSecretsMarshaller *SealedSecretMarshaller) Marshal(config core.Config, items []interface{}) (bool, []byte, error) {
	if !contains(config.OutputFormats, types.SEALED_SECRET) {
		return true, nil, nil
	}

	extraConfig, present := config.Extra[SEALED_SECRET_EXTRA_CONFIG_KEY]
	if !present {
		log.Fatalln("Unexpected missing kubeseal config")
	}
	kubesealConfig, ok := extraConfig.(KubesealConfig)
	if !ok {
		log.Fatalln("Unexpected format kubeseal config")
	}

	loadingRules := clientcmd.NewDefaultClientConfigLoadingRules()
	loadingRules.DefaultClientConfig = &clientcmd.DefaultClientConfig
	loadingRules.ExplicitPath = kubesealConfig.Kubeconfig
	var clientConfig kubeseal.ClientConfig = clientcmd.NewNonInteractiveDeferredLoadingClientConfig(loadingRules, &kubesealConfig.ConfigOverrides)

	pubKey, err := getPubKey(clientConfig, kubesealConfig)
	if err != nil {
		return false, nil, fmt.Errorf("failed retrieving public key for sealing secret: %v", err)
	}

	sealedSecrets := []ssv1.SealedSecret{}
	// do secret per secret because kubeseal doesn't support multiple secrets at once
	for _, item := range items {
		secret, ok := item.(corev1.Secret)
		if !ok {
			return false, nil, fmt.Errorf("unexpected format of item when unmarshalling to sealed secret: %#v", item)
		}
		bytes, err := marshalYaml([]corev1.Secret{secret})
		if err != nil {
			return false, nil, err
		}
		reader := b.NewReader(bytes)
		writer := new(b.Buffer)

		kubeseal.Seal(clientConfig, types.FormatIds[types.YAML][0], reader, writer, scheme.Codecs, pubKey, kubesealConfig.SealingScope,
			kubesealConfig.AllowEmptyData, "", "")

		sealed, err := io.ReadAll(writer)
		if err != nil {
			return false, nil, err
		}

		sealedSecret := ssv1.SealedSecret{}
		yaml.Unmarshal(sealed, &sealedSecret)
		sealedSecrets = append(sealedSecrets, sealedSecret)
	}

	res, err := marshal(config, sealedSecrets)
	if err != nil {
		return false, nil, err
	}
	return false, res, nil
}

func contains(formats []types.Format, str types.Format) bool {
	for _, v := range formats {
		if v == str {
			return true
		}
	}
	return false
}

func getPubKey(clientConfig kubeseal.ClientConfig, kubesealConfig KubesealConfig) (*rsa.PublicKey, error) {
	f, err := kubeseal.OpenCert(context.TODO(), clientConfig, kubesealConfig.ControllerNs, kubesealConfig.ControllerName, kubesealConfig.CertURL)
	if err != nil {
		return nil, err
	}
	defer f.Close()

	return kubeseal.ParseKey(f)
}

func marshal(config core.Config, sealedSecrets []ssv1.SealedSecret) ([]byte, error) {
	var outputFormat types.Format
	if len(config.OutputFormats) == 1 {
		outputFormat = config.InputFormat
	} else {
		for _, outputFormat = range config.OutputFormats {
			if outputFormat != types.SEALED_SECRET && (outputFormat == types.YAML || outputFormat == types.JSON) {
				break
			}
		}
	}
	if outputFormat == types.YAML {
		return marshalYaml(sealedSecrets)
	} else if outputFormat == types.JSON {
		return marshalJson(sealedSecrets)
	} else {
		return nil, fmt.Errorf("unable to find expected output format for sealed secret")
	}
}
