package cmd

import (
	"github.com/TwinProduction/gemplater/cmd/install"
	"github.com/spf13/cobra"
	"os"
)

func NewRootCmd() *cobra.Command {
	var rootCmd = &cobra.Command{
		Use:          "gemplater",
		Short:        "gemplater",
		Long:         "gemplater",
		Version:      "v0.0.1",
		SilenceUsage: true,
		Run: func(cmd *cobra.Command, args []string) {
			cmd.Help()
		},
	}
	return rootCmd
}

func Execute() {
	rootCmd := NewRootCmd()

	rootCmd.AddCommand(install.NewInstallCmd())

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
