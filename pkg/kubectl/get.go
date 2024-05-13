package kubectl

import "github.com/spf13/cobra"

var getCmd = &cobra.Command{
	Use:   "get <resources> | (<resource> <resource-name>)",
	Short: "get resources or get resource by resource name",
	Long:  "get resources or get resource by resource name",
	Args:  cobra.RangeArgs(1, 2),
	Run:   get,
}

func init() {
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "kube object' namespace")
	rootCmd.AddCommand(getCmd)
}
