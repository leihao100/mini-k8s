package kubectl

import (
	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "kubectl",
	Short: "cli for MiniK8S",
	Long:  "cli for MiniK8S",
}

func Execute() {
	rootCmd.PersistentFlags().StringP("filename", "f", "", "the name of yamlfile")
	rootCmd.PersistentFlags().StringP("namespace", "n", "", "kube object' namespace")
	rootCmd.AddCommand(applyCmd)
	rootCmd.AddCommand(getCmd)
	rootCmd.AddCommand(clearCmd)
	rootCmd.AddCommand(createCmd)
	rootCmd.AddCommand(deleteCmd)
	rootCmd.AddCommand(describeCmd)
	rootCmd.Execute()
}
