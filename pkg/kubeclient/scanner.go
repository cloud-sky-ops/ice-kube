package kubeclient

import (
	"context"
	"fmt"
	"log"
	"os"
	"path/filepath"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/tools/clientcmd"
)

// ScanCluster connects to the Kubernetes API and scans for unused resources
func ScanCluster(clusterName string) (string, error) {
	// Fetch kubeconfig file path to get config values like "KUBERNETES_SERVICE_HOST", "KUBERNETES_SERVICE_PORT" and many more
	// This logic counters the error ""unable to load in-cluster configuration, KUBERNETES_SERVICE_HOST and KUBERNETES_SERVICE_PORT must be defined""
	// The default path for kube config file is "/home/username/.kube/config" and this code extracts the same and stores all values in config varibale

	homeDir, err := os.UserHomeDir()
	if err != nil {
		log.Fatalf("Error getting home directory: %v", err)
	}

	kubeconfigPath := filepath.Join(homeDir, ".kube", "config")
	config, err := clientcmd.BuildConfigFromFlags("", kubeconfigPath)

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
