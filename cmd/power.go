package cmd

import (
	"fmt"

	"github.com/raws-labs/srig-cli/client"
	"github.com/raws-labs/srig-cli/output"
	"github.com/spf13/cobra"
)

func NewPowerCycleCmd(c **client.Client, jsonFlag *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "power-cycle",
		Short: "Power cycle the board in your active session",
		RunE: func(cmd *cobra.Command, args []string) error {
			sessionID, _ := cmd.Flags().GetString("session")

			var sess *client.Session
			var err error
			if sessionID != "" {
				sess, err = (*c).GetSession(sessionID)
				if err != nil {
					return fmt.Errorf("get session: %w", err)
				}
			} else {
				sess, err = (*c).FindActiveSession()
				if err != nil {
					return err
				}
			}

			if sess.State != "active" && sess.State != "idle" {
				return fmt.Errorf("session %s is %s, not active", sess.ID, sess.State)
			}

			if err := (*c).PowerCycle(sess.ID); err != nil {
				return err
			}

			if *jsonFlag {
				output.JSON(map[string]any{"success": true, "session_id": sess.ID})
			} else {
				output.Success(fmt.Sprintf("power cycle sent to session %s", sess.ID[:8]))
			}
			return nil
		},
	}

	cmd.Flags().String("session", "", "Session ID (auto-detects active session if omitted)")
	return cmd
}
