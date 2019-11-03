package install

import (
	"errors"
	"github.com/spf13/cobra"
	"io/ioutil"
	"os"
)

type Options struct {
	Destination string
}

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
				return errors.New("directory templating is not supported yet")
			} else {
				_, err = ioutil.ReadFile(fileName)
				if err != nil {
					return err
				}
			}
			return nil
		},
	}

	cmd.Flags().StringVarP(&options.Destination, "destination", "d", options.Destination, "Where to output the resulting file(s)")

	return cmd
}