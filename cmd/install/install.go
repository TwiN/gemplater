package install

import (
	"bufio"
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

	cfg, err := config.Get()
	// If the config hasn't been loaded, then load it
	if err == config.ErrConfigNotLoaded {
		if cfg, err = config.NewConfig(globalOptions.ConfigFile); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}

	cmd := &cobra.Command{
		Use:     "install FILE [DESTINATION]",
		Aliases: []string{"i"},
		Short:   "",
		Long:    "",
		Example: "gemplater install .profile ~/.profile",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]
			var destination string
			if len(args) > 1 {
				destination = args[1]
			}
			fi, err := os.Lstat(filePath)
			if err != nil {
				return err
			}
			fileOutputs := make(map[string]string)

			if fi.IsDir() {
				installDirectory(fileOutputs, filePath, cfg, options.IgnoreMissingVariables)
			} else {
				output, err := install(filePath, cfg, options.IgnoreMissingVariables)
				if err != nil {
					return err
				}
				fileOutputs[filePath] = output
			}

			for sourcePath, output := range fileOutputs {
				// If no destination provided, output to stdout
				if len(destination) == 0 {
					fmt.Printf("\n------ %s ------\n%s\n", sourcePath, output)
				} else {
					targetPath := strings.ReplaceAll(fmt.Sprintf("%s%s", destination, sourcePath[len(filePath):]), "\\", "/")
					elements := strings.Split(targetPath, "/")
					if len(elements) > 1 {
						targetParentPath := strings.Join(elements[:len(elements)-1], "/")
						if len(targetParentPath) != 0 {
							err = os.MkdirAll(targetParentPath, 0644)
							if err != nil {
								return err
							}
						}
					}
					fmt.Printf("Creating file at '%s' from template '%s'\n", targetPath, sourcePath)
					if err = ioutil.WriteFile(targetPath, []byte(output), 0644); err != nil {
						fmt.Printf("%v\n", err.Error())
						err = nil
					}
				}
			}
			return nil
		},
	}

	cmd.Flags().BoolVarP(&options.IgnoreMissingVariables, "ignore", "i", options.IgnoreMissingVariables, "Whether to ignore the missing variables")

	// TODO: --persist-choice
	// i.e. when interactiveChoice is called for file 1, the variables entered should be saved and could be reused
	// in file 2
	//

	return cmd
}

func installDirectory(fileOutputs map[string]string, filePath string, cfg *config.Config, ignoreMissingVariables bool) error {
	dir, err := ioutil.ReadDir(filePath)
	if err != nil {
		return err
	}
	for _, fi := range dir {
		path := fmt.Sprintf("%s%s%s", filePath, string(os.PathSeparator), fi.Name())
		if fi.IsDir() {
			installDirectory(fileOutputs, path, cfg, ignoreMissingVariables)
		} else {
			output, err := install(path, cfg, ignoreMissingVariables)
			if err != nil {
				fmt.Fprintf(os.Stderr, "%v", err.Error())
			}
			fileOutputs[path] = output
		}
	}
	return nil
}

func install(targetFilePath string, cfg *config.Config, ignoreMissingVariables bool) (string, error) {
	content, variables, err := processTargetFile(targetFilePath, cfg, ignoreMissingVariables)
	if err != nil {
		return "", err
	}
	return template.NewTemplate().WithVariables(variables).Replace(content), nil
}

func processTargetFile(targetFile string, cfg *config.Config, ignoreMissingVariables bool) (string, map[string]string, error) {
	rawContent, err := ioutil.ReadFile(targetFile)
	if err != nil {
		return "", nil, err
	}
	fileContent := string(rawContent)
	variables := cfg.Variables
	if !ignoreMissingVariables {
		err = interactiveVariables(targetFile, fileContent, variables)
		if err != nil {
			return "", nil, err
		}
	}
	return fileContent, variables, nil
}

func interactiveVariables(targetFile, fileContent string, variables map[string]string) error {
	printEvenIfSetInConfigFile := true // TODO: externalize that variable
	variableNames, err := ExtractVariablesFromString(fileContent, "__")
	if err != nil {
		return err
	}
	printedInstructions := false
	reader := bufio.NewReader(os.Stdin)
	for _, variableName := range variableNames {
		if _, exists := variables[variableName]; !exists || printEvenIfSetInConfigFile {
			if !printedInstructions {
				printedInstructions = true
				fmt.Printf("[%s]:\n", targetFile)
			}
			if exists && len(variables[variableName]) != 0 {
				fmt.Printf("Enter value for '%s' (default: %s): ", variableName, variables[variableName])
				value, _ := reader.ReadString('\n')
				value = strings.TrimSpace(value)
				if len(value) != 0 {
					variables[variableName] = value
				}
			} else {
				fmt.Printf("Enter value for '%s': ", variableName)
				value, _ := reader.ReadString('\n')
				variables[variableName] = strings.TrimSpace(value)
			}
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
