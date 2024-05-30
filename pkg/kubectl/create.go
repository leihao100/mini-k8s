package kubectl

import "github.com/spf13/cobra"

var createCmd = &cobra.Command{
	Use:   "create <resource> -f <filename>",
	Short: "create resource",
	Long:  "create resource",
	Args:  cobra.ExactArgs(2),
	Run:   create,
}
