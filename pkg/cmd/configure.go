package cmd

import (
	"bufio"
	"fmt"
	"os"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/virtualdom/tfdd/pkg/config"
)

type configureCmd struct {
	cfg *config.Config
}

// NewConfigureCmd returns a Cobra command that sets up the tfdd confuration.
func NewConfigureCmd() *cobra.Command {
	cc := &configureCmd{}

	cmd := &cobra.Command{
		Use:   "configure",
		Args:  cobra.ExactArgs(0),
		Short: "configure sets up AWS profile access",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			cc.cfg, err = getConfig(configFilePath)

			return err
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			var accountID string
			var name string
			var repeat string
			var err error

			r := bufio.NewReader(os.Stdin)

			for {
				fmt.Print("AWS Account ID: ")
				for {
					accountID, err = r.ReadString('\n')
					if err != nil {
						return errors.Wrap(err, "failed to read user input")
					}
					r.Reset(os.Stdin)

					accountID = strings.TrimSpace(accountID)
					_, err := strconv.ParseInt(accountID, 10, 64)
					if err == nil {
						break
					}

					fmt.Print("Please only input integers for the AWS Account ID: ")
				}

				fmt.Print("AWS Profile Name: ")
				name, err = r.ReadString('\n')
				if err != nil {
					return errors.Wrap(err, "failed to read user input")
				}
				r.Reset(os.Stdin)

				name = strings.TrimSpace(name)

				profile := &config.Profile{
					AccountID: accountID,
					Name:      name,
				}

				cc.cfg.Profiles = append(cc.cfg.Profiles, profile)

				fmt.Print("Add another profile? [Y/n]: ")
				for {
					repeat, err = r.ReadString('\n')
					if err != nil {
						return errors.Wrap(err, "failed to read user input")
					}
					r.Reset(os.Stdin)

					repeat = strings.TrimSpace(strings.ToLower(repeat))
					if repeat == "n" || repeat == "y" || repeat == "" {
						break
					}

					fmt.Print("Invalid input [Y/n]: ")
				}

				if repeat == "n" {
					break
				}
			}

			err = cc.cfg.Save()
			if err != nil {
				return err
			}

			fmt.Printf("Saved config to %v\n", cc.cfg.Path)
			return nil
		},
	}

	return cmd
}
