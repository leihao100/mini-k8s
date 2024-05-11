package kubectl

import "github.com/spf13/cobra"

var clearCmd = &cobra.Command{
	Use:   "clear",
	Short: "clear all resources",
	Long:  "clear all resources",
	Args:  cobra.MinimumNArgs(1),
	Run:   clear,
}

func init() {
	rootCmd.AddCommand(clearCmd)
}
