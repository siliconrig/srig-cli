package cmd

import (
	"fmt"

	"github.com/raws-labs/srig-cli/client"
	"github.com/raws-labs/srig-cli/output"
	"github.com/spf13/cobra"
)

func NewWhoamiCmd(c **client.Client, jsonFlag *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "whoami",
		Short: "Show your account info",
		RunE: func(cmd *cobra.Command, args []string) error {
			user, err := (*c).GetMe()
			if err != nil {
				return err
			}

			if *jsonFlag {
				output.JSON(user)
				return nil
			}

			output.Card("", [][2]string{
				{"email", user.Email},
				{"credits", fmt.Sprintf("%d", user.BalanceCredits)},
			})
			return nil
		},
	}
}
