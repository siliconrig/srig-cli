package main

import (
	"fmt"
	"os"

	"github.com/flashbay-dev/fbay-cli/client"
	"github.com/flashbay-dev/fbay-cli/cmd"
	"github.com/flashbay-dev/fbay-cli/output"
	"github.com/spf13/cobra"
)

var Version = "dev"

var (
	flagAPIKey  string
	flagBaseURL string
	flagJSON    bool
	c           *client.Client
)

var rootCmd = &cobra.Command{
	Use:   "fbay",
	Short: "flashbay — remote access to real MCU hardware",
	Long:  "CLI for the flashbay Hardware-as-a-Service platform.\nFlash firmware, open serial terminals, and run CI/CD tests on real boards.",
	PersistentPreRunE: func(cmd *cobra.Command, args []string) error {
		if flagAPIKey == "" {
			flagAPIKey = os.Getenv("FLASHBAY_API_KEY")
		}
		// Public commands — skip auth requirement.
		if cmd.Name() == "status" || cmd.Name() == "version" {
			c = client.New(flagBaseURL, flagAPIKey)
			return nil
		}
		if flagAPIKey == "" {
			return fmt.Errorf("no API key — set FLASHBAY_API_KEY or use --api-key")
		}
		c = client.New(flagBaseURL, flagAPIKey)
		return nil
	},
	SilenceUsage:  true,
	SilenceErrors: true,
}

func init() {
	rootCmd.PersistentFlags().StringVar(&flagAPIKey, "api-key", "", "API key (env: FLASHBAY_API_KEY)")
	rootCmd.PersistentFlags().StringVar(&flagBaseURL, "base-url", envOr("FLASHBAY_BASE_URL", "https://api.fbay.io"), "API base URL (env: FLASHBAY_BASE_URL)")
	rootCmd.PersistentFlags().BoolVar(&flagJSON, "json", false, "Output as JSON")

	rootCmd.AddCommand(cmd.NewStatusCmd(&c, &flagJSON))
	rootCmd.AddCommand(cmd.NewSessionCmd(&c, &flagJSON, &flagBaseURL))
	rootCmd.AddCommand(cmd.NewFlashCmd(&c, &flagJSON))
	rootCmd.AddCommand(cmd.NewSerialCmd(&c, &flagJSON))
	rootCmd.AddCommand(cmd.NewPowerCycleCmd(&c, &flagJSON))
	rootCmd.AddCommand(cmd.NewWhoamiCmd(&c, &flagJSON))
	rootCmd.AddCommand(&cobra.Command{
		Use:   "version",
		Short: "Print the CLI version",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("fbay " + Version)
		},
	})

	rootCmd.CompletionOptions.DisableDefaultCmd = true
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		output.Error(err.Error())
		os.Exit(1)
	}
}

func envOr(key, fallback string) string {
	if v := os.Getenv(key); v != "" {
		return v
	}
	return fallback
}
