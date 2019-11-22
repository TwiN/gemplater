package initialize

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/TwinProduction/gemplater/config"
	"github.com/TwinProduction/gemplater/core"
	"github.com/TwinProduction/gemplater/util"
	"github.com/spf13/cobra"
	"gopkg.in/yaml.v2"
	"io/ioutil"
	"os"
	"strings"
)

var (
	ErrFileNotFound = errors.New("cannot initialize a file that doesn't exist")
)

func NewInitializeCmd(globalOptions *core.GlobalOptions) *cobra.Command {
	cmd := &cobra.Command{
		Use:     "init [TARGET]",
		Aliases: []string{"I"},
		Short:   "Create a .gemplater.yml file with all variables found on target",
		Long:    "Create a .gemplater.yml file with all variables found on target",
		Example: "gemplater init\ngemplater init .profile",
		Args:    cobra.RangeArgs(0, 1),
		RunE: func(cmd *cobra.Command, args []string) error {
			if len(args) == 0 {
				args = append(args, ".")
			}
			target := args[0]
			fi, err := os.Lstat(target)
			if err != nil {
				if os.IsNotExist(err) {
					return ErrFileNotFound
				}
				return err
			}
			cfg := config.NewConfig(make(map[string]string))
			if fi.IsDir() {
				err = initializeDirectory(target, cfg)
			} else {
				err = initialize(target, cfg)
			}
			if err != nil {
				return err
			}
			configFileBytes, err := yaml.Marshal(cfg)
			if err != nil {
				return err
			}
			return ioutil.WriteFile(globalOptions.ConfigFile, configFileBytes, 0644)
		},
	}
	return cmd
}

func initializeDirectory(targetFile string, cfg *config.Config) error {
	dir, err := ioutil.ReadDir(targetFile)
	if err != nil {
		return err
	}
	for _, fi := range dir {
		path := fmt.Sprintf("%s%s%s", targetFile, string(os.PathSeparator), fi.Name())
		if fi.IsDir() {
			err = initializeDirectory(path, cfg)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%v", err.Error())
			}
		} else {
			err := initialize(path, cfg)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%v", err.Error())
			}
		}
	}
	return nil
}

func initialize(targetFile string, cfg *config.Config) error {
	targetFileBytes, err := ioutil.ReadFile(targetFile)
	if err != nil {
		return err
	}
	content := string(targetFileBytes)
	variableNames, err := util.ExtractVariablesFromString(content, "__")
	if err != nil {
		return err
	}
	printedFileName := false
	reader := bufio.NewReader(os.Stdin)
	for _, variableName := range variableNames {
		if !printedFileName {
			printedFileName = true
			fmt.Printf("\n[%s]:\n", targetFile)
		}
		if value, exists := cfg.Variables[variableName]; !exists {
			fmt.Printf("Enter value for '%s': ", variableName)
			value, _ := reader.ReadString('\n')
			cfg.Variables[variableName] = strings.TrimSpace(value)
		} else {
			fmt.Printf("Skipping variable '%s', because it is already set to '%s'\n", variableName, value)
			continue
		}
	}
	return nil
}
