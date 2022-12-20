package show

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
)

var (
	config = *core.NewConfig()
)

func init() {
	Cmd.PersistentFlags().StringSliceVarP(&config.InputFiles, "filename", "f", nil, "The files that contain the secrets to generate. Use '-' to read from stdin.")
	Cmd.MarkPersistentFlagFilename("filename")
	Cmd.PersistentFlags().VarP(enumflag.NewSlice(&config.OutputFormats, "output", types.FormatIds, enumflag.EnumCaseInsensitive), "output", "o", "Output format. One of: (json, yaml). If not specified the format of the input is used.")

	fs := Cmd.Flags()
	fs.AddGoFlagSet(goflag.CommandLine)
}

var Cmd = &cobra.Command{
	Use:     "show",
	Aliases: []string{"s"},
	Short:   "Show the template for a Secret that was previously create with kube-create-secret.",
	Long:    `Show the template for a Secret that was previously create with kube-create-secret.`,
	Example: "  kube-create-secret show -f secret.yml\n" +
		"  kube-create-secret show -f secret.json\n" +
		"  kube-create-secret show -f secret1.yml -f secret2.yml\n" +
		"  cat secret.yml | kube-create-secret show -f -\n" +
		"  kube-create-secret re-create -f secret.yml -o json\n",
	RunE: func(cmd *cobra.Command, args []string) error {
		if len(config.InputFiles) == 0 {
			return fmt.Errorf("reguired flag \"filename\" not set")
		}

		core.InputReaders = []core.InputReader{
			read.NewStdInReader(),
			read.NewFileReader(),
		}
		core.Unmarshallers = []core.Unmarshaller{
			unmarshal.NewReCreateJsonUnmarshaller(),
			unmarshal.NewReCreateYamlUnmarshaller(),
			unmarshal.NewReCreateJsonSealedSecretUnmarshaller(),
			unmarshal.NewReCreateYamlSealedSecretUnmarshaller(),
		}
		core.Marshallers = []core.Marshaller{
			marshal.NewYamlMarshaller(),
			marshal.NewJsonMarshaller(),
		}
		core.OutputWriters = []core.OutputWriter{
			write.NewStdOutWriter(),
		}

		config.InputReader = os.Stdin

		return core.ShowTemplate(config)
	},
}
