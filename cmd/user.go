package cmd

import (
	"fmt"
	"strings"

	client_query "github.com/jpeterburs/ts3-cli/internal"
	"github.com/spf13/cobra"
)

var (
	nickname     string
	muteInput    bool
	unmuteInput  bool
	muteOutput   bool
	unmuteOutput bool
	away         bool
	back         bool
	message      string
)

var userCmd = &cobra.Command{
	Use:   "user",
	Short: "Update your client's personal settings",
	Long:  "Update your client's personal settings, such as changing your nickname or toggling your input and output device",
	Example: `  Set your nickname:
    ts3 user --nickname 'John TeamSpeak'

  Mute yourself (input and output):
    ts3 user --mute-input --mute-output

  Set yourself as away with a message:
    ts3 user --away -m brb`,
	RunE: func(cmd *cobra.Command, args []string) error {
		client, err := client_query.Dial()
		if err != nil {
			return err
		}
		defer client.Quit()

		client.Authenticate()

		query := []string{"clientupdate"}

		if len(nickname) > 0 {
			query = append(query, fmt.Sprintf("client_nickname=%v", nickname))
		}
		if muteInput {
			query = append(query, "client_input_muted=1")
		}
		if unmuteInput {
			query = append(query, "client_input_muted=0")
		}
		if muteOutput {
			query = append(query, "client_output_muted=1")
		}
		if unmuteOutput {
			query = append(query, "client_output_muted=0")
		}
		if away {
			query = append(query, "client_away=1")
		}
		if back {
			query = append(query, "client_away=0")
		}
		if len(message) > 0 {
			query = append(query, fmt.Sprintf("client_away_message=%v", strings.ReplaceAll(message, " ", "\\s")))
		}

		client.Do(strings.Join(query, " "))

		cmd.Println("done")

		return nil
	},
}

func init() {
	rootCmd.AddCommand(userCmd)

	userCmd.Flags().StringVarP(&nickname, "nickname", "n", "", "Set a new nickname")

	userCmd.Flags().BoolVar(&muteInput, "mute-input", false, "Mute your input device")
	userCmd.Flags().BoolVar(&unmuteInput, "unmute-input", false, "Unmute your input device")
	userCmd.MarkFlagsMutuallyExclusive("mute-input", "unmute-input")

	userCmd.Flags().BoolVar(&muteOutput, "mute-output", false, "Mute your output device")
	userCmd.Flags().BoolVar(&unmuteOutput, "unmute-output", false, "Unmute your output device")
	userCmd.MarkFlagsMutuallyExclusive("mute-output", "unmute-output")

	userCmd.Flags().BoolVar(&away, "away", false, "Set away status")
	userCmd.Flags().BoolVar(&back, "back", false, "Unset your away status")
	userCmd.MarkFlagsMutuallyExclusive("away", "back")

	userCmd.Flags().StringVarP(&message, "message", "m", "", "Set what message to display when away")
}
