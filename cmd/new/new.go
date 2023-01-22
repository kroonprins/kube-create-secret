package new

import (
	goflag "flag"

	"github.com/bitnami-labs/sealed-secrets/pkg/pflagenv"
	"github.com/kroonprins/kube-create-secret/cmd/constants"
	"github.com/kroonprins/kube-create-secret/pkg/core"
	"github.com/kroonprins/kube-create-secret/pkg/output/marshal"
	"github.com/kroonprins/kube-create-secret/pkg/output/write"
	"github.com/kroonprins/kube-create-secret/pkg/types"
	"github.com/spf13/cobra"
	"github.com/thediveo/enumflag/v2"
)

var (
	config       = *core.NewConfig()
	templateType types.StarterTemplateType
)

func init() {
	fs := Cmd.PersistentFlags()
	fs.VarP(enumflag.NewSlice(&config.OutputFormats, "output", types.FormatIds, enumflag.EnumCaseInsensitive), "output", "o", "Output format. One of: (json, yaml). If not specified the format of the input is used.")
	fs.VarP(enumflag.New(&templateType, "type", types.StarterTemplateTypes, enumflag.EnumCaseInsensitive), "type", "t", "Template type. One of: (data, stringData, tls).")

	fs.AddGoFlagSet(goflag.CommandLine)
	pflagenv.SetFlagsFromEnv(constants.FLAGENV_PREFIX, fs)
}

var Cmd = &cobra.Command{
	Use:     "new",
	Aliases: []string{"n"},
	Short:   "Print starter template.",
	Long:    `Print starter template.`,
	Example: "  kube-create-secret new\n" +
		"  kube-create-secret new -t tls -o json\n",
	RunE: func(cmd *cobra.Command, args []string) error {
		core.Marshallers = []core.Marshaller{
			marshal.NewYamlMarshaller(),
			marshal.NewJsonMarshaller(),
		}
		core.OutputWriters = []core.OutputWriter{
			write.NewStdOutWriter(),
		}

		return core.StarterTemplate(config, templateType)
	},
}
