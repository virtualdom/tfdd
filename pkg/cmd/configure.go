package cmd

import (
	"bufio"
	"fmt"
	"os"
	"path"
	"path/filepath"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/virtualdom/tfdd/pkg/config"
)

type configureCmd struct {
	configFilePath string

	cfg *config.Config
}

func NewConfigureCmd() *cobra.Command {
	cc := &configureCmd{}

	cmd := &cobra.Command{
		Use:   "configure",
		Args:  cobra.ExactArgs(0),
		Short: "configure sets up AWS profile access",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			if len(cc.configFilePath) == 0 {
				cc.configFilePath = path.Join(os.Getenv("HOME"), ".tfdd", "config")
			}

			cc.configFilePath, err = filepath.Abs(cc.configFilePath)
			if err != nil {
				return errors.Wrap(err, fmt.Sprintf("failed to get absolute path of %v", cc.configFilePath))
			}

			cc.cfg, err = config.New(cc.configFilePath)
			return errors.Wrap(err, fmt.Sprintf("failed to load config file %v", cc.configFilePath))
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
					Name: name,
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
						break;
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

			fmt.Printf("Saved config to %v", cc.cfg.Path)
			return nil
		},
	}

	f := cmd.Flags()
	f.StringVarP(&cc.configFilePath, "config-file", "c", os.Getenv("TFDD_CONFIG_FILE"), "specify config file. Can also set environment variable `TFDD_CONFIG_FILE`. Defaults to `~/.tfdd/config`.")

	return cmd
}
