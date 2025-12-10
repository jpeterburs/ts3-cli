package cmd

import (
	"fmt"
	"strings"

	client_query "github.com/jpeterburs/ts3-cli/internal"
	"github.com/spf13/cobra"
)

var selfCmd = &cobra.Command{
	Use:                   "self [-n nickname] [--mute-input | --unmute-input] [--mute-output | --unmute-output] [--away | --back] [-m message]",
	DisableFlagsInUseLine: true,
	Short:                 "Update your client's personal settings",
	Long:                  "Update your client's personal settings, such as changing your nickname or toggling your input and output device",
	Example: `Set your nickname:
  ts3 self --nickname 'John TeamSpeak'

Mute yourself (input and output):
  ts3 self --mute-input --mute-output

Set yourself as away with a message:
  ts3 self --away -m brb`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := client_query.Dial()
		if err != nil {
			return err
		}
		defer client.Quit()

		if err := client.Authenticate(); err != nil {
			return err
		}

		query := []string{"clientupdate"}

		nickname, _ := cmd.Flags().GetString("nickname")
		if len(nickname) > 0 {
			query = append(query, fmt.Sprintf("client_nickname=%v", strings.ReplaceAll(nickname, " ", "\\s")))
		}

		muteInput, _ := cmd.Flags().GetBool("mute-input")
		if muteInput {
			query = append(query, "client_input_muted=1")
		}

		unmuteInput, _ := cmd.Flags().GetBool("unmute-input")
		if unmuteInput {
			query = append(query, "client_input_muted=0")
		}

		muteOutput, _ := cmd.Flags().GetBool("mute-output")
		if muteOutput {
			query = append(query, "client_output_muted=1")
		}

		unmuteOutput, _ := cmd.Flags().GetBool("unmute-output")
		if unmuteOutput {
			query = append(query, "client_output_muted=0")
		}

		away, _ := cmd.Flags().GetBool("away")
		if away {
			query = append(query, "client_away=1")
		}

		back, _ := cmd.Flags().GetBool("back")
		if back {
			query = append(query, "client_away=0")
		}

		message, _ := cmd.Flags().GetString("message")
		if len(message) > 0 {
			query = append(query, fmt.Sprintf("client_away_message=%v", strings.ReplaceAll(message, " ", "\\s")))
		}

		if err := client.Do(strings.Join(query, " ")); err != nil {
			return err
		}

		cmd.Println("done")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(selfCmd)

	selfCmd.Flags().StringP("nickname", "n", "", "Set a new nickname")

	selfCmd.Flags().Bool("mute-input", false, "Mute your input device")
	selfCmd.Flags().Bool("unmute-input", false, "Unmute your input device")
	selfCmd.MarkFlagsMutuallyExclusive("mute-input", "unmute-input")

	selfCmd.Flags().Bool("mute-output", false, "Mute your output device")
	selfCmd.Flags().Bool("unmute-output", false, "Unmute your output device")
	selfCmd.MarkFlagsMutuallyExclusive("mute-output", "unmute-output")

	selfCmd.Flags().Bool("away", false, "Set away status")
	selfCmd.Flags().Bool("back", false, "Unset your away status")
	selfCmd.MarkFlagsMutuallyExclusive("away", "back")

	selfCmd.Flags().StringP("message", "m", "", "Set what message to display when away")
}
