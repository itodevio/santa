package santa

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

func RequireConfig(_ *cobra.Command, _ []string) {
	if !viper.IsSet("SessionToken") {
		if err := viper.ReadInConfig(); err != nil {
			fmt.Println("Config not found. Please run `santa config --session <session-token>` to set your token globally.")
			os.Exit(0)
		}
	} else {
		fmt.Printf("SessionToken: %s\n", viper.GetString("SessionToken"))
	}
}
