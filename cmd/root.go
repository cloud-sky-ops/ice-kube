package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var clusterName string
var deleteBeforeHours int
var namespace string
var dryRun bool

// rootCmd represents the base command
var rootCmd = &cobra.Command{
	Use:   "ice-kube",
	Short: "CLI tool for cleaning up unused Kubernetes resources",
	Long: `ICE Kube is a command-line tool to help identify and remove 
unused Kubernetes resources for better cost optimization.`,
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println("ICE Kube CLI - Use 'ice-kube scan' to scan for unused resources.")
	},
}

// Execute adds all child commands to the root command and sets flags appropriately.
func Execute() {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Println("Error executing command:", err)
	}
}

func init() {
	// Place to define global flags and configuration settings.
	// Global flags are supported by cobra using PersistentFlags
	rootCmd.PersistentFlags().BoolP("verbose", "v", false, "Share detailed logs")
}
