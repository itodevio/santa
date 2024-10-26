package cmd

import (
	"fmt"
	"os"
	"path"
	"path/filepath"

	"github.com/adrg/xdg"
	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var rootCmd = &cobra.Command{
	Use:   "santa",
	Short: "Santa is an utility tool for Advent of Code problems",
	Long:  "Santa is an utility tool for Advent of Code problems built with love by itodevio in Go",
	PersistentPreRun: func(cmd *cobra.Command, args []string) {
		err := viper.ReadInConfig()
		if err != nil && cmd.Name() != "config" {
			fmt.Println("No config file found. Run `santa config` to create one.")
			os.Exit(0)
		}
	},
}

func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("An error ocurred: %v\n", err)
	}
}

func init() {
	cobra.OnInitialize(initConfig)
}

func initConfig() {
	home, err := os.UserHomeDir()
	if err != nil {
		panic(err)
	}

	viper.SetConfigType("yaml")

	matches, _ := filepath.Glob(".santa*")
	if len(matches) > 0 {
		viper.SetConfigName(".santa")
		viper.AddConfigPath(".")
	} else {
		viper.SetConfigName("config")
		viper.AddConfigPath(path.Join(home, ".santa"))
		viper.AddConfigPath(path.Join(xdg.ConfigHome, "santa"))
	}
}
