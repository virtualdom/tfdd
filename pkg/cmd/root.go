package cmd

import (
	"fmt"
	"os"
	"path/filepath"

	"github.com/spf13/cobra"
)

var (
	rootDir string
	version string // version can be overwritten by ldflags in the Makefile.
)

// rootCmd represents the base command when called without any subcommands.
var rootCmd = &cobra.Command{
	Use:   "tfdd",
	Args:  cobra.ExactArgs(1),
	Short: "tfdd explores a Terraform project and checks for drift",
	Long: `tfdd (Terraform Drift Detector) takes a Terraform project root as
input and visits all subdirectories to check for drift. Created
by Github user @virtualdom in Go.
Complete documentation is available at https://github.com/virtualdom/tfdd`,
	Version: version,

	PreRun: func(cmd *cobra.Command, args []string) {
		path, err := filepath.Abs(args[0])
		if err != nil {
			fmt.Println("failed to get absolute path of " + args[0])
			os.Exit(1)
		}

		if _, err = os.Stat(path); os.IsNotExist(err) {
			fmt.Println(path + " not found")
			os.Exit(1)
		}

		rootDir = path
	},

	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(rootDir)
	},
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
	rootCmd.Flags().BoolP("version", "v", false, "print Captain version number")

	// Add Captain commands.
	// rootCmd.AddCommand(command goes here)

	// Prevent usage message from being printed out upon command error.
	rootCmd.SilenceUsage = true
}
