package cmd

import (
	"context"
	"fmt"

	"github.com/google/go-github/github"
	version "github.com/hashicorp/go-version"
	"github.com/spf13/cobra"
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

func upgradeRun(cmd *cobra.Command, args []string) {
	latest, err := getLatestVersionFromGithub()
	if err != nil {
		fmt.Println("Failed to get latest version from Github:", err)
		return
	}

	current := Version
	fmt.Println("Latest:", latest, "Current:", current)
}

func init() {
	rootCmd.AddCommand(upgradeCmd)
	upgradeCmd.Flags().BoolVarP(&yes, "yes", "y", false, "Automatically answer 'yes' to any prompt that might appear on the command line")
}
