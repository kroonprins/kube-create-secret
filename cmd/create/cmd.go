package create

import (
	goflag "flag"
	"fmt"
	"os"

	"github.com/kroonprins/kube-create-secret/pkg/core"
	"github.com/kroonprins/kube-create-secret/pkg/input/read"
	"github.com/kroonprins/kube-create-secret/pkg/input/unmarshal"
	"github.com/kroonprins/kube-create-secret/pkg/output/marshal"
	"github.com/kroonprins/kube-create-secret/pkg/output/write"
	"github.com/kroonprins/kube-create-secret/pkg/types"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/klog/v2"
)

var (
	config          = *core.NewConfig()
	kubesealConfig  marshal.KubesealConfig
	configOverrides clientcmd.ConfigOverrides
)

func init() {
	Cmd.PersistentFlags().StringSliceVarP(&config.InputFiles, "filename", "f", nil, "The files that contain the secrets to generate. Use '-' to read from stdin.")
	Cmd.MarkPersistentFlagFilename("filename")
	Cmd.PersistentFlags().VarP(enumflag.NewSlice(&config.OutputFormats, "output", types.FormatIds, enumflag.EnumCaseInsensitive), "output", "o", "Output format. One of: (json, yaml, sealed-secret). If not specified the format of the input is used. In case of sealed-secret this can be combined with a second format json or yaml to specifiy the format of the SealedSecret")

	Cmd.PersistentFlags().StringVar(&kubesealConfig.CertURL, "kubeseal-cert", "", "Only relevant if output format is sealed-secret. Certificate / public key file/URL to use for encryption. Overrides --controller-*")
	Cmd.PersistentFlags().StringVar(&kubesealConfig.ControllerNs, "kubeseal-controller-namespace", metav1.NamespaceSystem, "Only relevant if output format is sealed-secret. Namespace of sealed-secrets controller.")
	Cmd.PersistentFlags().StringVar(&kubesealConfig.ControllerName, "kubeseal-controller-name", "sealed-secrets-controller", "Only relevant if output format is sealed-secret. Name of sealed-secrets controller.")
	Cmd.PersistentFlags().BoolVar(&kubesealConfig.AllowEmptyData, "kubeseal-allow-empty-data", false, "Only relevant if output format is sealed-secret. Allow empty data in the secret object")
	Cmd.PersistentFlags().StringVar(&kubesealConfig.Kubeconfig, "kubeseal-kubeconfig", "", "Only relevant if output format is sealed-secret. Path to a kube config. Only required if out-of-cluster")
	Cmd.PersistentFlags().Var(&kubesealConfig.SealingScope, "kubeseal-scope", "Only relevant if output format is sealed-secret. Set the scope of the sealed secret: strict, namespace-wide, cluster-wide (defaults to strict). Mandatory for --raw, otherwise the 'sealedsecrets.bitnami.com/cluster-wide' and 'sealedsecrets.bitnami.com/namespace-wide' annotations on the input secret can be used to select the scope.")

	klog.InitFlags(nil)

	fs := Cmd.Flags()
	fs.AddGoFlagSet(goflag.CommandLine)

	kflags := clientcmd.RecommendedConfigOverrideFlags("kubeseal-k8s-")
	clientcmd.BindOverrideFlags(&configOverrides, fs, kflags)
}

var Cmd = &cobra.Command{
	Use:     "create",
	Aliases: []string{"c"},
	Short:   "Create a secret from a SecretTemplate definition.",
	Long:    `Create a secret from a SecretTemplate definition.`,
	Example: "  kube-create-secret create -f template.yml\n" +
		"  kube-create-secret create -f template.json\n" +
		"  kube-create-secret create -f template1.yml -f template2.yml\n" +
		"  cat template.yml | kube-create-secret create -f -\n" +
		"  kube-create-secret create -f template.yml -o json\n" +
		"  kube-create-secret create -f template.yml -o json -o sealed-secret\n",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(config.InputFiles) == 0 {
			return fmt.Errorf("reguired flag \"filename\" not set")
		}

		core.InputReaders = []core.InputReader{
			read.NewStdInReader(),
			read.NewFileReader(),
		}
		core.Unmarshallers = []core.Unmarshaller{
			unmarshal.NewJsonUnmarshaller(),
			unmarshal.NewYamlUnmarshaller(), // json is put before yaml because yaml unmarshaller succeeds for json input
		}
		core.Marshallers = []core.Marshaller{
			marshal.NewYamlMarshaller(),
			marshal.NewJsonMarshaller(),
			marshal.NewSealedSecretMarshaller(),
		}
		core.OutputWriters = []core.OutputWriter{
			write.NewStdOutWriter(),
		}

		config.InputReader = os.Stdin
		kubesealConfig.ConfigOverrides = configOverrides
		config.Extra[marshal.SEALED_SECRET_EXTRA_CONFIG_KEY] = kubesealConfig

		return core.Create(config)
	},
}
