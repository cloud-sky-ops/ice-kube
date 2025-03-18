package cmd

import (
	"fmt"
	"os"

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
		result, err := kubeclient.ScanCluster(clusterName)
		if err != nil {
			fmt.Println("Error scanning cluster:", err)
			os.Exit(1)
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
}
