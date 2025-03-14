package kubeclient

import (
	"fmt"
	"context"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// ScanCluster connects to the Kubernetes API and scans for unused resources
func ScanCluster(clusterName string) (string, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return "", fmt.Errorf("failed to create in-cluster config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	// Example: Fetching all pods in the cluster
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list pods: %v", err)
	}

	return fmt.Sprintf("Found %d pods in cluster %s", len(pods.Items), clusterName), nil
}
