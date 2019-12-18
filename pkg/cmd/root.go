package cmd

import (
	"os"

	"github.com/spf13/cobra"
)

var (
	configFilePath string
	version        string // version can be overwritten by ldflags in the Makefile.
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "tfdd",
	Short: "tfdd explores a Terraform project and checks for drift",
	Long: `tfdd (Terraform Drift Detector) takes a Terraform project root as
input and visits all subdirectories to check for drift. Run "tfdd configure"
to begin. Created by Github user @virtualdom in Go.
Complete documentation is available at https://github.com/virtualdom/tfdd`,
	Version: version,
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	return rootCmd.Execute()
}

func init() {
	// Set default version number if none given.
	if rootCmd.Version == "" {
		rootCmd.Version = "0.0.0"
	}

	// Add flags.
	rootCmd.Flags().BoolP("version", "v", false, "print tfdd version number")
	rootCmd.PersistentFlags().StringVarP(&configFilePath, "config-file", "c", os.Getenv("TFDD_CONFIG_FILE"), "specify config file. Can also set environment variable `TFDD_CONFIG_FILE`. Defaults to `~/.tfdd/config`.")

	// Add tfdd commands.
	rootCmd.AddCommand(NewConfigureCmd())
	rootCmd.AddCommand(NewDetectCmd())

	// Prevent usage message from being printed out upon command error.
	rootCmd.SilenceUsage = true
}
