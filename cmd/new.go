package cmd

import (
	"embed"
	"errors"
	"fmt"
	"io"
	"net/http"
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
				return errors.New("too many arguments")
			}
			if !isValidYear(viper.GetInt("year")) {
				return errors.New("year config not valid. run `santa config --year x` to set a valid aoc year")
			}

			now := time.Now()
			month := 12
			day, err := getDay(args)
			fmt.Println(viper.GetInt("year"), month, day)
			if err != nil {
				return errors.New("invalid day argument")
			}

			aocDate := time.Date(viper.GetInt("year"), time.Month(month), day, 0, 0, 0, 0, aocLocation)
			if now.Before(aocDate) {
				return errors.New("day not yet available. specify a different day in the argument or practice in the previous year")
			}

			return nil
		},
	}
	//go:embed templates/*
	templatesFS embed.FS
)

func newRun(cmd *cobra.Command, args []string) {
	day, err := getDay(args)
	if err != nil {
		panic(err)
	}

	dayPath := fmt.Sprintf("Day%s", formatDay(day))
	err = os.Mkdir(dayPath, os.ModePerm)
	if err != nil { // TODO: better error handling here
		panic(err)
	}

	mainTmpl, err := template.ParseFS(templatesFS, "templates/go/main.tmpl")
	if err != nil { // TODO: better error handling here
		panic(err)
	}
	part1Tmpl, err := template.ParseFS(templatesFS, "templates/go/part1.tmpl")
	if err != nil { // TODO: better error handling here
		panic(err)
	}
	part2Tmpl, err := template.ParseFS(templatesFS, "templates/go/part2.tmpl")
	if err != nil { // TODO: better error handling here
		panic(err)
	}

	downloadInput(day, dayPath)
	createTemplateFile(mainTmpl, dayPath, "main.go")
	createTemplateFile(part1Tmpl, dayPath, "part1.go")
	createTemplateFile(part2Tmpl, dayPath, "part2.go")
}

func createTemplateFile(tmpl *template.Template, dayPath string, filename string) {
	file, err := os.Create(path.Join(dayPath, filename))
	if err != nil { // TODO: better error handling here
		panic(err)
	}
	defer file.Close()

	err = tmpl.Execute(file, nil)
	if err != nil {
		panic(err)
	}
}

func downloadInput(day int, dayPath string) {
	url := fmt.Sprintf("https://adventofcode.com/%d/day/%d/input", viper.GetInt("year"), day)
	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to create request:", err)
		return
	}

	sessionCookie := &http.Cookie{
		Name:  "session",
		Value: viper.GetString("session"),
	}

	req.AddCookie(sessionCookie)

	client := &http.Client{}
	resp, err := client.Do(req)
	if err != nil {
		fmt.Println("Failed to send request:", err)
		return
	}
	defer resp.Body.Close()

	input, err := io.ReadAll(resp.Body)
	if err != nil {
		fmt.Println("Failed to read response:", err)
		return
	}

	file, err := os.Create(fmt.Sprintf("%s/input.txt", dayPath))
	if err != nil {
		fmt.Println("Failed to create input file:", err)
		return
	}
	defer file.Close()

	file.Write(input)
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
