package initialize

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/TwinProduction/gemplater/config"
	"github.com/TwinProduction/gemplater/core"
	"github.com/TwinProduction/gemplater/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

var (
	ErrFileNotFound = errors.New("cannot initialize a file that doesn't exist")
)

// Create config file with all __VAR_NAMES__ in current folder and allows you to set the default value
// (perhaps this shouldn't be a subcommand of config?)

func NewInitializeCmd(globalOptions *core.GlobalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init [TARGET]",
		Aliases: []string{"I"},
		Short:   "Create .gemplater.yml file with all variables found on target",
		Long:    "Create .gemplater.yml file with all variables found on target",
		Example: "gemplater init\ngemplater init .profile",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			// TODO: if len(args) == 0...
			target := args[0]
			fi, err := os.Lstat(target)
			if err != nil {
				if os.IsNotExist(err) {
					return ErrFileNotFound
				}
				return err
			}

			if fi.IsDir() {
				return errors.New("initializing directory isn't supported yet")
			} else {
				initialize(target, globalOptions.ConfigFile)
			}
			return nil
		},
	}
	return cmd
}

func initialize(targetFile, configFile string) error {
	raw, err := ioutil.ReadFile(targetFile)
	if err != nil {
		return err
	}
	content := string(raw)
	variableNames, err := util.ExtractVariablesFromString(content, "__")
	if err != nil {
		return err
	}
	cfg := config.NewConfig(make(map[string]string))
	printedInstructions := false
	reader := bufio.NewReader(os.Stdin)
	for _, variableName := range variableNames {
		if value, exists := cfg.Variables[variableName]; !exists {
			if !printedInstructions {
				printedInstructions = true
				fmt.Printf("[%s]:\n", targetFile)
			}
			fmt.Printf("Enter value for '%s': ", variableName)
			value, _ := reader.ReadString('\n')
			cfg.Variables[variableName] = strings.TrimSpace(value)
		} else {
			fmt.Printf("Skipping variable '%s', because it is already set to '%s'", variableName, value)
			continue
		}
	}

	return nil
}
