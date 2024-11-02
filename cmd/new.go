package cmd

import (
	"errors"
	"fmt"
	"os"
	"path"
	"strconv"
	"text/template"
	"time"

	"github.com/spf13/cobra"
	"github.com/spf13/viper"
)

var (
	newCmd = &cobra.Command{
		Use:   "new [day]",
		Short: "Creates new boilerplate for selected day, or today if in december and no day specified.",
		Run:   newRun,
		Args: func(cmd *cobra.Command, args []string) error {
			aocLocation, err := time.LoadLocation("America/New_York")
			if err != nil {
				panic(err)
			}
			if len(args) > 1 {
				return errors.New("Too many arguments.")
			}
			if !isValidYear(viper.GetInt("year")) {
				return errors.New("Year config not valid. Run `santa config --year x` to set a valid AoC year.")
			}

			now := time.Now()
			month := 12
			day, err := getDay(args)
			fmt.Println(viper.GetInt("year"), month, day)
			if err != nil {
				return errors.New("Invalid day argument.")
			}

			aocDate := time.Date(viper.GetInt("year"), time.Month(month), day, 0, 0, 0, 0, aocLocation)
			if now.Before(aocDate) {
				return errors.New("Day not yet available. Specify a different day in the argument or practice in the previous year!")
			}

			return nil
		},
	}
)

func newRun(cmd *cobra.Command, args []string) {
	day, err := getDay(args)
	if err != nil {
		panic(err)
	}

	dayDir := fmt.Sprintf("Day%s", formatDay(day))
	err = os.Mkdir(dayDir, os.ModePerm)
	if err != nil { // TODO: better error handling here
		panic(err)
	}

	tmpl, err := template.ParseFiles("cmd/templates/go.tmpl")
	if err != nil { // TODO: better error handling here
		panic(err)
	}

	part1, err := os.Create(path.Join(dayDir, "solution1.go"))
	if err != nil { // TODO: better error handling here
		panic(err)
	}
	defer part1.Close()

	part2, err := os.Create(path.Join(dayDir, "solution2.go"))
	if err != nil { // TODO: better error handling here
		panic(err)
	}
	defer part2.Close()

	err = tmpl.Execute(part1, nil)
	if err != nil {
		panic(err)
	}
	err = tmpl.Execute(part2, nil)
	if err != nil {
		panic(err)
	}

	// req, err := http.Get("")
}

func formatDay(day int) string {
	strDay := strconv.Itoa(day)

	if len(strDay) == 1 {
		return "0" + strDay
	}

	return strDay
}

func getDay(args []string) (int, error) {
	var err error
	day := time.Now().Day()
	if len(args) > 0 {
		day, err = strconv.Atoi(args[0])
		if err != nil {
			return 0, err
		}
	}

	return day, nil
}

func init() {
	rootCmd.AddCommand(newCmd)

}
