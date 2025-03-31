package cmd

import (
	"fmt"

	utils "github.com/cloud-sky-ops/ice-kube/internal"
	"github.com/cloud-sky-ops/ice-kube/pkg/kubeclient"
	"github.com/spf13/cobra"
)

var scanCmd = &cobra.Command{
	Use:   "scan",
	Short: "Scan a Kubernetes cluster for unused resources",
	RunE: func(cmd *cobra.Command, args []string) error {
		if clusterName == "" {
			return fmt.Errorf("cluster name is required. Use --cluster or -c flag")
		}

		fmt.Println("Scanning cluster:", clusterName)
		result, err := kubeclient.ScanCluster(clusterName, deleteBeforeHours)
		if err != nil {
			message := "Error scanning cluster:"
			utils.PrintError(message, err)
		}
		fmt.Println("Scan Result:", result)
		return nil
	},
}

func init() {
	// Add scanCmd to the root command
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "Specify the cluster name")
	scanCmd.MarkFlagRequired("cluster")
	scanCmd.Flags().IntVarP(&deleteBeforeHours, "delete-before-hours", "", 24, "Delete resources created before these number of hours. Default is 24 hours." )
}
