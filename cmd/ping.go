package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var pingCmd = &cobra.Command{
	Use:   "ping",
	Short: "Play ping pong with Santa! He's a master at it.",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("Pong!")
		fmt.Println(viper.AllSettings())
	},
}

func init() {
	rootCmd.AddCommand(pingCmd)
}
