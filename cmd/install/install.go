package install

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/TwinProduction/gemplater/config"
	"github.com/TwinProduction/gemplater/core"
	"github.com/TwinProduction/gemplater/template"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

type Options struct {
	IgnoreMissingVariables bool
}

func NewInstallCmd(globalOptions *core.GlobalOptions) *cobra.Command {
	options := &Options{}

	cmd := &cobra.Command{
		Use:     "install FILE [DESTINATION]",
		Aliases: []string{"i"},
		Short:   "",
		Long:    "",
		Example: "gemplater install .profile ~/.profile",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileName := args[0]
			var destination string
			if len(args) > 1 {
				destination = args[1]
			}
			fi, err := os.Lstat(fileName)
			if err != nil {
				return err
			}
			if fi.IsDir() {
				// Read all files one by one
				// ...
				//if len(options.Destination) == 0 {
				//	fmt.Printf("%s/%s\n%s\n\n", destination, fileName, output)
				//}
				return errors.New("directory templating is not supported yet")
			} else {
				rawContent, err := ioutil.ReadFile(fileName)
				if err != nil {
					return err
				}
				content := string(rawContent)
				variables, err := getVariables(globalOptions.ConfigFile, content, options.IgnoreMissingVariables)
				if err != nil {
					return err
				}
				output := template.NewTemplate().WithVariables(variables).Replace(content)
				// If no destination provided, the output will be stdout
				if len(destination) == 0 {
					println(output)
				} else {
					fmt.Printf("Create file at '%s' from template '%s'\n", destination, fileName)
					return ioutil.WriteFile(destination, []byte(output), 0644)
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&options.IgnoreMissingVariables, "ignore", "i", options.IgnoreMissingVariables, "Whether to ignore the missing variables")

	return cmd
}

func getVariables(configFile, templateFileContent string, ignoreMissingVariables bool) (variables map[string]string, err error) {
	cfg, err := config.Get()
	// If the config hasn't been loaded, then load it
	if err == config.ErrConfigNotLoaded {
		if cfg, err = config.NewConfig(configFile); err != nil {
			return nil, err
		}
	} else if err != nil {
		return nil, err
	}
	variables = cfg.Variables
	if !ignoreMissingVariables {
		err = interactiveVariables(templateFileContent, variables)
		if err != nil {
			return nil, err
		}
	}
	return
}

func interactiveVariables(fileContent string, variables map[string]string) error {
	variableNames, err := ExtractVariablesFromString(fileContent, "__")
	if err != nil {
		return err
	}
	reader := bufio.NewReader(os.Stdin)
	for _, variable := range variableNames {
		if _, exists := variables[variable]; !exists {
			fmt.Printf("Enter value for '%s': ", variable)
			value, _ := reader.ReadString('\n')
			variables[variable] = strings.TrimSpace(value)
		}
	}
	return nil
}

func ExtractVariablesFromString(s, wrapper string) (variableNames []string, err error) {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	// Instead of doing it all at once, we'll do it line by line to reduce the odds of picking up a multiline variable
	for _, line := range lines {
		variablesInLine := getAllBetween(line, wrapper, wrapper)
		for _, variable := range variablesInLine {
			if strings.Contains(variable, " ") {
				continue
			}
			variableNames = append(variableNames, variable)
		}
	}
	return
}

// Get all substrings between two strings
// This variation does not strip the suffix and the prefix from the substring
func getAllBetween(haystack, prefix, suffix string) (needles []string) {
	for {
		if len(haystack) < len(prefix)+len(suffix) {
			break
		}
		start := strings.Index(haystack, prefix) + len(prefix)
		if start-len(prefix) == -1 {
			break
		}
		end := strings.Index(haystack[start:], suffix) + start
		if end-start == -1 || start >= end {
			break
		}
		needles = append(needles, haystack[start-len(prefix):end+len(suffix)])
		if len(haystack) <= end {
			break
		}
		haystack = haystack[end+len(suffix):]
	}
	return needles
}
