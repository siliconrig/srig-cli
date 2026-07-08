package cmd

import (
	"context"
	"encoding/json"
	"fmt"
	"os"
	"regexp"
	"strconv"
	"time"

	"github.com/siliconrig/srig-cli/client"
	"github.com/siliconrig/srig-cli/firmware"
	"github.com/siliconrig/srig-cli/output"
	"github.com/spf13/cobra"
)

func NewFlashCmd(c **client.Client, jsonFlag *bool) *cobra.Command {
	cmd := &cobra.Command{
		Use:   "flash <firmware>",
		Short: "Flash firmware to a remote board",
		Long: "Upload and flash firmware to the board in your active session.\n" +
			"Accepts a raw .bin (all boards), a .uf2 (rp2350), or an .elf / Intel .hex\n" +
			"for STM32 boards (auto-converted to a raw image before upload).\n" +
			"If --session is not given, auto-detects the active session.",
		Example: "  srig flash build/app.bin\n" +
			"  srig flash build/firmware.elf --session sess_abc123\n" +
			"  srig flash build/firmware.hex",
		Args: cobra.ExactArgs(1),
		RunE: func(cmd *cobra.Command, args []string) error {
			filePath := args[0]
			sessionID, _ := cmd.Flags().GetString("session")
			timeout, _ := cmd.Flags().GetDuration("timeout")

			data, err := os.ReadFile(filePath)
			if err != nil {
				return fmt.Errorf("cannot read firmware file: %w", err)
			}

			var sess *client.Session
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

			flashBytes, fwInfo, err := firmware.Normalize(data, sess.BoardType)
			if err != nil {
				return err
			}
			if len(flashBytes) > 16<<20 {
				return fmt.Errorf("firmware is %d bytes, over the platform's 16 MB limit", len(flashBytes))
			}
			if !*jsonFlag {
				if fwInfo.Format != firmware.FormatRaw {
					output.Info(fmt.Sprintf("detected %s, converted to %d-byte bin @ %#x", fwInfo.Format, fwInfo.Size, fwInfo.BaseAddr))
				}
				output.Info(fmt.Sprintf("flashing %s (%d bytes) to session %s", filePath, len(flashBytes), sess.ID))
			}

			ctx, cancel := context.WithTimeout(context.Background(), timeout)
			defer cancel()

			ws, err := (*c).DialSerialWS(ctx, sess.ID)
			if err != nil {
				return fmt.Errorf("connect to session: %w", err)
			}
			defer ws.CloseNow()

			if err := (*c).FlashFirmwareBytes(sess.ID, "firmware.bin", flashBytes); err != nil {
				return fmt.Errorf("upload firmware: %w", err)
			}

			if !*jsonFlag {
				output.Info("upload complete, flashing...")
			}

			// Wait for flash_done via WS
			pctRe := regexp.MustCompile(`Writing at.*?(\d+(?:\.\d+)?)\s*%`)
			for {
				_, data, err := ws.Read(ctx)
				if err != nil {
					return fmt.Errorf("lost connection while flashing: %w", err)
				}

				var msg map[string]any
				if json.Unmarshal(data, &msg) != nil {
					continue
				}

				switch msg["type"] {
				case "flash_output":
					line, _ := msg["line"].(string)
					if *jsonFlag {
						output.JSON(map[string]string{"type": "progress", "line": line})
					} else {
						if m := pctRe.FindStringSubmatch(line); m != nil {
							pct, _ := strconv.ParseFloat(m[1], 64)
							fmt.Printf("\r  %s", output.ProgressBar(pct, 30))
						}
					}
				case "flash_done":
					success, _ := msg["success"].(bool)
					errMsg, _ := msg["error"].(string)

					if !*jsonFlag {
						fmt.Println()
					}

					if success {
						if *jsonFlag {
							output.JSON(map[string]any{"success": true})
						} else {
							output.Success("flash complete")
						}
						return nil
					}
					if *jsonFlag {
						output.JSON(map[string]any{"success": false, "error": errMsg})
					}
					return fmt.Errorf("flash failed: %s", errMsg)
				}
			}
		},
	}

	cmd.Flags().String("session", "", "Session ID (auto-detects active session if omitted)")
	cmd.Flags().Duration("timeout", 2*time.Minute, "Timeout for flash operation")
	return cmd
}
