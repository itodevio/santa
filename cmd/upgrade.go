package cmd

import (
	"bufio"
	"context"
	"fmt"
	"io"
	"net/http"
	"os"
	"path/filepath"
	"runtime"
	"slices"
	"strings"

	"github.com/google/go-github/github"
	version "github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
	"golang.org/x/text/cases"
	"golang.org/x/text/language"
)

var (
	upgradeCmd = &cobra.Command{
		Use:   "upgrade",
		Short: "Upgrade Santa to the latest version",
		Run:   upgradeRun,
	}
	yes bool
)

func getLatestVersionFromGithub() (*version.Version, error) {
	gh := github.NewClient(nil)
	ctx := context.Background()
	tags, _, err := gh.Repositories.ListTags(ctx, "itodevio", "santa", nil)
	if err != nil {
		return nil, err
	}

	var versionsRaw []string
	for _, tag := range tags {
		versionsRaw = append(versionsRaw, *tag.Name)
	}

	latest := version.Must(version.NewVersion("0.0.0"))
	for _, versionRaw := range versionsRaw {
		v, err := version.NewVersion(versionRaw)
		if err != nil {
			return nil, err
		}

		if v.GreaterThan(latest) {
			latest = v
		}
	}

	return latest, nil
}

func getARMVersion() (string, error) {
	data, err := os.ReadFile("/proc/cpuinfo")
	if err != nil {
		return "", err
	}

	info := string(data)
	switch {
	case strings.Contains(info, "ARMv5"):
		return "armv5l", nil
	case strings.Contains(info, "ARMv6"):
		return "armv6l", nil
	case strings.Contains(info, "ARMv7"):
		return "armv7l", nil
	case strings.Contains(info, "ARMv8"):
		return "armv8l", nil
	}

	return "arm", nil
}

func getDownloadURL(version string) (string, error) {
	var err error
	supportedByOS := map[string][]string{
		"darwin":  {"arm64", "x86_64"},
		"linux":   {"arm", "arm64", "armv5l", "armv6l", "armv7l", "armv8l", "i386", "x86_64"},
		"windows": {"i386", "amd64", "arm", "arm64"},
	}

	os := runtime.GOOS
	arch := runtime.GOARCH
	ext := ""

	if os == "linux" && arch == "arm" {
		arch, err = getARMVersion()
		if err != nil {
			return "", err
		}
	}

	switch arch {
	case "386":
		arch = "i386"
	case "amd64":
		arch = "x86_64"
	}

	if os == "windows" {
		ext = ".exe"
	}

	if !slices.Contains(supportedByOS[os], arch) {
		return "", fmt.Errorf("unsupported OS/ARCH: %s/%s", os, arch)
	}

	caser := cases.Title(language.English)

	return fmt.Sprintf(
		"https://github.com/itodevio/santa/releases/download/v%s/santa-%s-%s%s",
		version,
		caser.String(os),
		arch,
		ext,
	), nil
}

func upgradeRun(cmd *cobra.Command, args []string) {
	latest, err := getLatestVersionFromGithub()
	if err != nil {
		fmt.Println("Failed to get latest version from Github:", err)
		return
	}

	current := version.Must(version.NewVersion(Version))
	if !current.LessThan(latest) {
		fmt.Println("You are already using the latest version!")
		return
	}

	fmt.Printf("There's a new version available.\nCurrent: %s\nLatest: %s\n", current.String(), latest.String())

	if !yes {
		fmt.Print("Do you want to upgrade? (Y/n) ")
		response, err := bufio.NewReader(os.Stdin).ReadString('\n')
		if err != nil {
			fmt.Println("Failed to read response:", err)
			return
		}

		response = strings.ToLower(strings.TrimSpace(response))
		if response == "n" || response == "no" {
			return
		}
	}

	selfPath, err := os.Executable()
	if err != nil {
		fmt.Println("Failed to get executable path:", err)
		return
	}

	selfDir := filepath.Dir(selfPath)

	file, err := os.Open(selfPath)
	if err != nil {
		fmt.Println("Failed to open executable file:", err)
		return
	}
	defer file.Close()

	stat, err := file.Stat()
	if err != nil {
		fmt.Println("Failed to get file info:", err)
		return
	}

	newFileName := "temp_" + filepath.Base(selfPath)
	newFilePath := filepath.Join(selfDir, newFileName)

	newFile, err := os.Create(newFilePath)
	if err != nil {
		if strings.Contains(strings.ToLower(err.Error()), "permission denied") {
			fmt.Println("Permission denied. Try running the command with 'sudo'.")
			return
		}
		fmt.Println("Failed to create new file:", err)
		return
	}
	defer newFile.Close()
	defer os.Remove(newFilePath)

	err = newFile.Chmod(stat.Mode())
	if err != nil {
		fmt.Println("Failed to change new file permissions:", err)
		return
	}

	url, err := getDownloadURL(latest.String())
	if err != nil {
		fmt.Println("Failed to upgrade santa:", err)
		return
	}

	request, err := http.NewRequest("GET", url, nil)
	if err != nil {
		fmt.Println("Failed to upgrade santa:", err)
		return
	}

	client := &http.Client{}
	response, err := client.Do(request)
	if response.StatusCode != 200 {
		err = fmt.Errorf("failed to download binary from GitHub. HTTP Status: %d", response.StatusCode)
	}
	if err != nil {
		fmt.Println("Failed to upgrade santa:", err)
		return
	}
	defer response.Body.Close()

	_, err = io.Copy(newFile, response.Body)
	if err != nil {
		fmt.Println("Failed to save new binary:", err)
		return
	}

	err = os.Rename(newFilePath, selfPath)
	if err != nil {
		fmt.Println("Failed to replace current binary:", err)
		return
	}

	fmt.Println("Santa has been upgraded to the latest version!")
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	upgradeCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatically answer 'yes' to any prompt that might appear on the command line")
}
