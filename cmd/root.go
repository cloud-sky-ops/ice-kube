package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

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
