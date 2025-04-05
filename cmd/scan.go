package cmd

import (
	"fmt"

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

		if namespace != "" {
			fmt.Printf("Scanning cluster: %s in namespace %s\n", clusterName, namespace)
		} else {
			fmt.Println("Scanning resources in cluster:", clusterName)
		}
		pods, _,  err := kubeclient.ScanResources(clusterName, namespace)
		if err != nil {
			fmt.Println("Error in scanning resources in cluster:", clusterName)
			return err
		}
		fmt.Printf("Found %d pods in cluster %s\n", len(pods.Items), clusterName)
		return nil
	},
}

func init() {
	rootCmd.AddCommand(scanCmd)
	scanCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "Specify the cluster name")
	scanCmd.MarkFlagRequired("cluster")
	scanCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace to scan pods. Default is all namespaces")
}
