package commands

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	cfgFile     string
	verbose     bool
	quiet       bool
	noColor     bool
	showVersion bool
)

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "easyClean",
	Short: "Detect and remove unused assets from your codebase",
	Long: `easyClean is a CLI tool that automatically detects and safely removes
unused asset files (images, fonts, videos, etc.) from codebases.

It uses smart scanning with multi-pattern reference detection and supports
multiple project types (React, Vue, Flutter, iOS, Android).`,
	Version: "1.0.0",
}

// Execute runs the root command
func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}

func init() {
	// Global flags
	rootCmd.PersistentFlags().StringVarP(&cfgFile, "config", "c", ".unusedassets.yaml", "config file path")
	rootCmd.PersistentFlags().BoolVarP(&verbose, "verbose", "v", false, "enable verbose logging")
	rootCmd.PersistentFlags().BoolVarP(&quiet, "quiet", "q", false, "suppress all output except errors")
	rootCmd.PersistentFlags().BoolVar(&noColor, "no-color", false, "disable colored output")
}

// GetConfigFile returns the config file path
func GetConfigFile() string {
	return cfgFile
}

// IsVerbose returns whether verbose mode is enabled
func IsVerbose() bool {
	return verbose
}

// IsQuiet returns whether quiet mode is enabled
func IsQuiet() bool {
	return quiet
}

// IsColorDisabled returns whether color output is disabled
func IsColorDisabled() bool {
	return noColor
}
