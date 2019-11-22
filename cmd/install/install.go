package install

import (
	"bufio"
	"fmt"
	"github.com/TwinProduction/gemplater/config"
	"github.com/TwinProduction/gemplater/core"
	"github.com/TwinProduction/gemplater/template"
	"github.com/TwinProduction/gemplater/util"
	"github.com/spf13/cobra"
	"io/ioutil"
	"net/http"
	"os"
	"strings"
)

type Options struct {
	IgnoreMissingVariables bool // -i --ignore
	Quick                  bool // -q --quick
	Remember               bool // -r --remember
}

func NewInstallCmd(globalOptions *core.GlobalOptions) *cobra.Command {
	options := &Options{
		IgnoreMissingVariables: false,
		Quick:                  false,
		Remember:               false,
	}
	cfg, err := config.Get()
	// If the config hasn't been loaded, then load it
	if err == config.ErrConfigNotLoaded {
		if cfg, err = config.ReadConfig(globalOptions.ConfigFile); err != nil {
			panic(err)
		}
	} else if err != nil {
		panic(err)
	}
	cmd := &cobra.Command{
		Use:     "install TARGET [DESTINATION]",
		Aliases: []string{"i"},
		Short:   "",
		Long:    "",
		Example: "gemplater install .profile ~/.profile\ngemplater install .profile --quick --remember",
		Args:    cobra.RangeArgs(1, 2),
		RunE: func(cmd *cobra.Command, args []string) error {
			target := args[0]
			var destination string
			if len(args) > 1 {
				destination = args[1]
			}
			fi, err := os.Lstat(target)
			if err != nil {
				return err
			}
			fileOutputs := make(map[string]string)
			// TODO: support URL target
			if fi.IsDir() {
				installDirectory(fileOutputs, target, cfg, options)
			} else {
				output, err := install(target, cfg, options)
				if err != nil {
					return err
				}
				fileOutputs[target] = output
			}
			for sourcePath, output := range fileOutputs {
				// If no destination provided, output to stdout
				if len(destination) == 0 {
					fmt.Printf("\n------ %s ------\n%s\n", sourcePath, output)
				} else {
					targetPath := strings.ReplaceAll(fmt.Sprintf("%s%s", destination, sourcePath[len(target):]), "\\", "/")
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

	cmd.Flags().BoolVarP(&options.IgnoreMissingVariables, "ignore", "i", options.IgnoreMissingVariables, "Whether to ignore missing variables. If not set, missing variables will trigger interactive mode")
	cmd.Flags().BoolVarP(&options.Quick, "quick", "q", options.Quick, "Do not ask for value of variables that are already set. Requires -i to not be set")
	cmd.Flags().BoolVarP(&options.Remember, "remember", "r", options.Remember, "Remember variables interactively set on one file for other files. Requires -i to not be set. Useless if TARGET is not directory")

	return cmd
}

func installDirectory(fileOutputs map[string]string, filePath string, cfg *config.Config, options *Options) error {
	dir, err := ioutil.ReadDir(filePath)
	if err != nil {
		return err
	}
	for _, fi := range dir {
		path := fmt.Sprintf("%s%s%s", filePath, string(os.PathSeparator), fi.Name())
		if fi.IsDir() {
			installDirectory(fileOutputs, path, cfg, options)
		} else {
			output, err := install(path, cfg, options)
			if err != nil {
				_, _ = fmt.Fprintf(os.Stderr, "%v", err.Error())
			}
			fileOutputs[path] = output
		}
	}
	return nil
}

func install(targetFilePath string, cfg *config.Config, options *Options) (string, error) {
	content, variables, err := processTargetFile(targetFilePath, cfg, options)
	if err != nil {
		return "", err
	}
	return template.NewTemplate().WithVariables(variables).Replace(content), nil
}

func processTargetFile(targetFile string, cfg *config.Config, options *Options) (content string, variables map[string]string, err error) {
	raw, err := ioutil.ReadFile(targetFile)
	if err != nil {
		return "", nil, err
	}
	content = string(raw)
	if options.Remember {
		variables = cfg.Variables
	} else {
		variables = make(map[string]string)
		for k, v := range cfg.Variables {
			variables[k] = v
		}
	}
	if !options.IgnoreMissingVariables {
		err = interactiveVariables(targetFile, content, variables, options)
		if err != nil {
			return "", nil, err
		}
	}
	return
}

func interactiveVariables(targetFile, fileContent string, variables map[string]string, options *Options) error {
	variableNames, err := util.ExtractVariablesFromString(fileContent, "__")
	if err != nil {
		return err
	}
	printedInstructions := false
	reader := bufio.NewReader(os.Stdin)
	for _, variableName := range variableNames {
		if defaultValue, exists := variables[variableName]; !exists || options.Quick {
			if !printedInstructions {
				printedInstructions = true
				fmt.Printf("[%s]:\n", targetFile)
			}
			if exists && len(variables[variableName]) != 0 {
				fmt.Printf("Enter value for '%s' (default: %s): ", variableName, defaultValue)
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

func getFileContentFromUrl(url string) (string, error) {
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	data, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return "", err
	}
	return string(data), nil
}
