package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	initCmd = &cobra.Command{
		Use:   "init [project_path]",
		Args:  validateArgs,
		Short: "Init's an Advent of Code project folder",
		Run:   initRun,
	}
	defaultPath   = "."
	ProjectPath   = defaultPath
	forceNonEmpty = false
)

func initRun(cmd *cobra.Command, args []string) {
	err := os.Mkdir(ProjectPath, os.ModePerm)
	if err != nil && !os.IsExist(err) {
		fmt.Println("Failed to create project directory:", err)
		return
	}

	viper.WriteConfigAs(path.Join(ProjectPath, ".santa.yaml"))
	fmt.Printf("Project created at %s, cd into it to create the first day!", ProjectPath)
}

func validateArgs(cmd *cobra.Command, args []string) error {
	if len(args) == 0 {
		return nil
	}
	if len(args) > 1 {
		return errors.New(fmt.Sprintf("Command init expected at most 1 argument and %d were given.", len(args)))
	}

	path := args[0]
	info, err := os.Stat(path)
	if os.IsNotExist(err) {
		ProjectPath = path
		return nil
	}
	if err != nil {
		return errors.New("Failed to read path. " + err.Error())
	}
	if !info.IsDir() {
		return errors.New("Target path is not a directory.")
	}

	files, err := os.ReadDir(path)
	if err != nil {
		return errors.New("Failed to read path. " + err.Error())
	}
	if len(files) != 0 && !forceNonEmpty {
		return errors.New("Directory is not empty. Run command with --force to ignore warning.")
	}

	ProjectPath = path
	return nil
}

func isValidYear(year int) bool {
	now := time.Now()
	aoCFirstYear := 2015

	if year < now.Year() && year >= aoCFirstYear {
		return true
	}
	if year == now.Year() && now.Month() == 12 {
		return true
	}

	return false
}

func getLatestAoCYear() int {
	now := time.Now()
	if now.Month() == 12 {
		return now.Year()
	}

	return now.Year() - 1
}

func init() {
	rootCmd.AddCommand(initCmd)
	initCmd.Flags().Int("year", getLatestAoCYear(), "The year of the AoC you want to do. Defaults to current year if in december or later, or to previous year otherwise.")
	initCmd.Flags().BoolVar(&forceNonEmpty, "force", false, "Whether should force the project init in a non empty directory.")
	viper.BindPFlag("year", initCmd.Flags().Lookup("year"))
}
