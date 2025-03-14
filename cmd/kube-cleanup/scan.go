package kube_cleanup

import (
	"fmt"
	"os"

	"github.com/cloud-sky-ops/kube-cleanup/pkg/kubeclient"

	"github.com/spf13/cobra"
)

func main() {
	var clusterName string

	var rootCmd = &cobra.Command{
		Use:   "kube-cleanup",
		Short: "A CLI tool for cleaning up unused Kubernetes resources",
	}

	var scanCmd = &cobra.Command{
		Use:   "scan",
		Short: "Scan a Kubernetes cluster for unused resources",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("Scanning cluster:", clusterName)
			result, err := kubeclient.ScanCluster(clusterName)
			if err != nil {
				fmt.Println("Error scanning cluster:", err)
				os.Exit(1)
			}
			fmt.Println("Scan Result:", result)
		},
	}

	scanCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "Specify the cluster name")
	scanCmd.MarkFlagRequired("cluster")
	rootCmd.AddCommand(scanCmd)

	if err := rootCmd.Execute(); err != nil {
		fmt.Println(err)
		os.Exit(1)
	}
}
