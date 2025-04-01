package cmd

import (
	"fmt"

	utils "github.com/cloud-sky-ops/ice-kube/internal"
	"github.com/cloud-sky-ops/ice-kube/pkg/kubeclient"
	"github.com/spf13/cobra"
)

var deleteCmd = &cobra.Command{
	Use:   "delete",
	Short: "Delete unused resources from a Kubernetes cluster",
	RunE: func(cmd *cobra.Command, args []string) error {
		if clusterName == "" {
			return fmt.Errorf("cluster name is required. Use --cluster or -c flag")
		}

		if namespace != "" {
			fmt.Printf("Deleting cluster: %s in namespace %s\n", clusterName, namespace)
		} else {
			fmt.Println("Deleting resources in cluster:", clusterName)
		}
		result, err := kubeclient.DeleteResources(clusterName, deleteBeforeHours, namespace)
		if err != nil {
			message := "Error deleting resources:"
			utils.PrintError(message, err)
		}
		fmt.Println("Delete Result:", result)
		return nil
	},
}

func init() {
	// Add deleteCmd to the root command
	rootCmd.AddCommand(deleteCmd)
	deleteCmd.Flags().StringVarP(&clusterName, "cluster", "c", "", "Specify the cluster name")
	deleteCmd.MarkFlagRequired("cluster")
	deleteCmd.Flags().IntVarP(&deleteBeforeHours, "delete-before-hours", "t", 24, "Delete resources created before these number of hours.")
	deleteCmd.Flags().StringVarP(&namespace, "namespace", "n", "", "Namespace to delete resources. Default is all namespaces")
}
