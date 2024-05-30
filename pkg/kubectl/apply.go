package kubectl

import (
	"github.com/spf13/cobra"
)

var applyCmd = &cobra.Command{
	Use:   "apply -f <filename>",
	Short: "apply resource",
	Long:  "apply resource",
	Run:   create,
}
