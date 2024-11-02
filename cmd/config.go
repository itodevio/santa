package cmd

import (
	"fmt"
	"os"
	"path"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var configCmd = &cobra.Command{
	Use:   "config",
	Short: "Set up config options to a local config file or globally",
	Run:   configRun,
}

func configRun(cmd *cobra.Command, args []string) {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	err = viper.ReadInConfig()
	if err != nil {
		configPath := path.Join(home, ".santa")
		configFile := path.Join(configPath, "config")
		err := os.MkdirAll(configPath, 0755)
		if err != nil {
			panic(err)
		}

		_, err = os.Create(configFile)
		if err != nil {
			panic(err)
		}

		err = viper.WriteConfigAs(configFile)
		if err != nil {
			panic(err)
		}

		fmt.Printf("Wrote config to '%s'.\n", configFile)
		return
	}

	err = viper.WriteConfig()
	if err != nil {
		panic(err)
	}

	fmt.Printf("Wrote config to '%s'.\n", viper.ConfigFileUsed())
}

func init() {
	rootCmd.AddCommand(configCmd)
	configCmd.Flags().StringP("session", "s", "", "The session token found in the browser's network tab in the input page")
	configCmd.MarkFlagRequired("session")
	configCmd.Flags().Int("year", getLatestAoCYear(), "The year of the AoC you want to do. Defaults to current year if in december or later, or to previous year otherwise.")
	viper.BindPFlag("session", configCmd.Flags().Lookup("session"))
	viper.BindPFlag("year", configCmd.Flags().Lookup("year"))
}
