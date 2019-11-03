package install

import (
	"bufio"
	"errors"
	"fmt"
	"github.com/TwinProduction/gemplater/template"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
	"strings"
)

type Options struct {
	Destination string
}

var (
	// Since the wrapper should be present at the beginning and the end of the variable name
	ErrOddNumberOfSeparator = errors.New("uneven number of separator")
)

func NewInstallCmd() *cobra.Command {
	options := &Options{}

	cmd := &cobra.Command{
		Use:     "install FILE",
		Aliases: []string{"i"},
		Short:   "",
		Long:    "",
		Example: "gemplater install .profile",
		Args:    cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			fileName := args[0]
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
				variables, err := InteractiveVariables(content)
				if err != nil {
					return err
				}
				_ = template.NewTemplate().WithVariables(variables).Replace(content)
				// If no destination provided, the output will be stdout
				if len(options.Destination) == 0 {
					//fmt.Println(output)
				}
			}
			return nil
		},
	}

	// This overrides the file configuration
	cmd.Flags().StringVarP(&options.Destination, "destination", "d", options.Destination, "Where to output the resulting file(s). If no value is specified, the output will be stdout")
	// TODO: flag to ignore missing variables

	return cmd
}

func InteractiveVariables(fileContent string) (map[string]string, error) {
	variableNames, err := ExtractVariablesFromString(fileContent, "__")
	if err != nil {
		return nil, err
	}
	variables := make(map[string]string)
	reader := bufio.NewReader(os.Stdin)
	for _, variable := range variableNames {
		fmt.Printf("Enter value for '%s': ", variable)
		value, _ := reader.ReadString('\n')
		variables[variable] = value
	}
	return variables, nil
}

func ExtractVariablesFromString(s, wrapper string) (variableNames []string, err error) {
	s = strings.ReplaceAll(s, "\r\n", "\n")
	lines := strings.Split(s, "\n")
	// Instead of doing it all at once, we'll do it line by line to reduce the odds of picking up a multiline variable
	for _, line := range lines {
		//elements := strings.Split(line, wrapper)
		////if len(elements) % 2 == 0 {
		////	continue
		////}
		//for i, element := range elements {
		//
		//	if i % 2 != 0 {
		//		variableNames = append(variableNames, fmt.Sprintf("%s%s%s", wrapper, element, wrapper))
		//		println(fmt.Sprintf("%s%s%s", wrapper, element, wrapper))
		//	}
		//}
		variable, _ := getBetween(line, wrapper, wrapper)
		if len(variable) > 0 {
			variableNames = append(variableNames, fmt.Sprintf("%s%s%s", wrapper, variable, wrapper))
			println(fmt.Sprintf("%s%s%s", wrapper, variable, wrapper))
		}
	}
	return
}

// Get substring between two strings.
func getBetween(value string, a string, b string) (string, int) {
	if len(value) < len(a)+len(b) {
		return "", -1
	}
	start := strings.Index(value, a) + len(a)
	if start == -1 {
		return "", -1
	}
	end := strings.Index(value[start:], b) + start
	if end == -1 {
		return "", -1
	}
	if start >= end {
		return "", -1
	}
	return value[start:end], end
}
