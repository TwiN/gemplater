package cmd

import (
	"github.com/TwinProduction/gemplater/cmd/initialize"
	"github.com/TwinProduction/gemplater/cmd/install"
	"github.com/TwinProduction/gemplater/core"
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
	globalOptions := core.NewGlobalOptions(".gemplater.yml")

	rootCmd.PersistentFlags().StringVarP(&globalOptions.ConfigFile, "config", "c", globalOptions.ConfigFile, "Specify configuration file to use")

	rootCmd.AddCommand(install.NewInstallCmd(globalOptions))
	rootCmd.AddCommand(initialize.NewInitializeCmd(globalOptions))

	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}
