package kubectl

import "github.com/spf13/cobra"

var createCmd = &cobra.Command{
	Use:   "create <resource> -f <filename>",
	Short: "create resource",
	Long:  "create resource",
	Args:  cobra.MinimumNArgs(1),
	Run:   create,
}

func init() {
	rootCmd.AddCommand(createCmd)
}
