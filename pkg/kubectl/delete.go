package kubectl

import "github.com/spf13/cobra"

var deleteCmd = &cobra.Command{
	Use:   "delete <resource> <resource-name>",
	Short: "delete resource",
	Long:  "delete resource",
	Args:  cobra.MinimumNArgs(1),
	Run:   delete,
}
