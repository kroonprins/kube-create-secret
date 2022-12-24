package main

import (
	"fmt"
	"os"

	"github.com/kroonprins/kube-create-secret/cmd/create"
	re_create "github.com/kroonprins/kube-create-secret/cmd/re-create"
	"github.com/kroonprins/kube-create-secret/cmd/show"
	"github.com/spf13/cobra"
)

var (
	Version = "0.0.0"
)

func newRootCmd() *cobra.Command {
	rootCmd := &cobra.Command{
		Use:   "kube-create-secret",
		Short: "Utility for creating kubernetes secrets.",
		Long:  `Utility for creating kubernetes secrets.`,
		Example: "  kube-create-secret create -f template.yaml\n" +
			"  kube-create-secret re-create -f secret.yaml\n" +
			"  kube-create-secret show -f secret.yaml\n",
	}

	versionCmd := &cobra.Command{
		Use:     "version",
		Short:   "Print the version",
		Example: "  kube-create-secret version\n",
		Long:    ``,
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Printf("Version: %s\n", Version)
		},
	}

	rootCmd.AddCommand(versionCmd, create.Cmd, re_create.Cmd, show.Cmd)

	return rootCmd
}

// func init() {
// 	klog.InitFlags(nil)
// }

func main() {
	rootCmd := newRootCmd()
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "FAILURE: '%s'\n", err)
		os.Exit(1)
	}
}
