package cmd

import (
	"os"
	"path/filepath"

	"github.com/aws/aws-sdk-go/service/sts"
	"github.com/pkg/errors"
	"github.com/spf13/cobra"
	"github.com/virtualdom/tfdd/pkg/auth"
	"github.com/virtualdom/tfdd/pkg/aws"
	"github.com/virtualdom/tfdd/pkg/config"
)

type detectCmd struct {
	cfg           *config.Config
	rootDir       string
	profile       string
	accessKey     string
	secretKey     string
	sessionToken  string
	auditRoleName string
	auth          bool
	authInterval  int
}

// NewDetectCmd returns a Cobra command that explores a Terraform project and
// generates Terraform plans to detect drift.
func NewDetectCmd() *cobra.Command {
	dc := &detectCmd{}

	cmd := &cobra.Command{
		Use:   "detect",
		Args:  cobra.ExactArgs(1),
		Short: "detect explores a Terraform project to detect drift",

		PreRunE: func(cmd *cobra.Command, args []string) error {
			var err error

			dc.cfg, err = getConfig(configFilePath)
			if err != nil {
				return err
			}

			path, err := filepath.Abs(args[0])
			if err != nil {
				return errors.Wrap(err, "failed to get absolute path of "+args[0])
			}

			if _, err = os.Stat(path); os.IsNotExist(err) {
				return errors.New(path + " not found")
			}

			dc.rootDir = path
			return nil
		},

		RunE: func(cmd *cobra.Command, args []string) error {
			if dc.auth {
				var sts *sts.STS
				var err error

				switch {
				case dc.profile != "":
					sts, err = aws.NewSTSClient(dc.profile)
				case dc.accessKey != "" && dc.secretKey != "" && dc.sessionToken != "":
					sts, err = aws.NewSTSClient(dc.accessKey, dc.secretKey, dc.sessionToken)
				case dc.accessKey != "" && dc.secretKey != "":
					sts, err = aws.NewSTSClient(dc.accessKey, dc.secretKey)
				default:
					sts, err = aws.NewSTSClient()
				}

				if err != nil {
					return errors.Wrap(err, "failed to initialize AWS client")
				}

				auth.Auth(dc.cfg, sts, dc.auditRoleName)

				if dc.authInterval >= 0 {
					go auth.Interval(dc.cfg, sts, dc.auditRoleName, dc.authInterval)
				}

				// continue detect code here! Just needed an easy entrypoint to test auth.
			}

			return nil
		},
	}

	cmd.Flags().BoolVarP(&dc.auth, "auth", "a", false, "specify whether to assume an AWS role (specified by --audit-role-name), and write credentials to ~/.aws/credentials. Use with --interval to periodically refresh credentials.")
	cmd.Flags().IntVarP(&dc.authInterval, "interval", "i", 0, "number of seconds between AWS role assumptions. If omitted, the role will be assumed once at the beginning and never again. This will only take effect if --auth is set.")
	cmd.Flags().StringVarP(&dc.auditRoleName, "audit-role-name", "r", "audit-role", "the name of the IAM role to assume in all AWS accounts.")
	cmd.Flags().StringVarP(&dc.profile, "profile", "p", "", "AWS profile to use when attempting role assumption. Falls back on default AWS SDK authentication behavior if omitted. Do not use with --access-key, --secret-key, or --session-token.")
	cmd.Flags().StringVar(&dc.accessKey, "access-key", "", "IAM access key to use when attempting role assumption. Falls back on default AWS SDK authentication behavior if omitted. Only use with keys and session tokens. Do not use with --profile.")
	cmd.Flags().StringVar(&dc.secretKey, "secret-key", "", "IAM secret key to use when attempting role assumption. Falls back on default AWS SDK authentication behavior if omitted. Only use with keys and session tokens. Do not use with --profile.")
	cmd.Flags().StringVar(&dc.sessionToken, "session-token", "", "IAM session token to use when attempting role assumption. Falls back on default AWS SDK authentication behavior if omitted. Only use with access/secret keys. Do not use with --profile.")

	return cmd
}
