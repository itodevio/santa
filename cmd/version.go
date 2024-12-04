package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

func init() {
	rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Print santa's version",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(Version)
	},
}
