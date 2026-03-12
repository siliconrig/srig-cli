package cmd

import (
	"fmt"

	"github.com/flashbay-dev/fbay-cli/client"
	"github.com/flashbay-dev/fbay-cli/output"
	"github.com/spf13/cobra"
)

func NewStatusCmd(c **client.Client, jsonFlag *bool) *cobra.Command {
	return &cobra.Command{
		Use:   "status",
		Short: "Show board availability",
		RunE: func(cmd *cobra.Command, args []string) error {
			boards, err := (*c).ListBoards()
			if err != nil {
				return err
			}

			if *jsonFlag {
				output.JSON(boards)
				return nil
			}

			if len(boards) == 0 {
				output.Success("no boards registered")
				return nil
			}

			avail := 0
			for _, b := range boards {
				if b.State == "available" {
					avail++
				}
			}

			headers := []string{"TYPE", "LABEL", "STATE", "ID"}
			rows := make([][]string, len(boards))
			for i, b := range boards {
				rows[i] = []string{b.BoardType, b.Label, b.State, b.ID}
			}

			fmt.Println()
			output.Table(headers, rows)
			fmt.Printf("\n  %s %s%s\n",
				output.StyleDim.Render("available:"),
				output.StyleGreen.Render(fmt.Sprintf("%d", avail)),
				output.StyleDim.Render(fmt.Sprintf("/%d", len(boards))))
			return nil
		},
	}
}
