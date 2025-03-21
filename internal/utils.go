package utils

import (
	"fmt"
	"log"
	"os"
	"path/filepath"

	"github.com/TwiN/go-color"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)
// Note: Only function names starting with capital letters are allowed to be 
// 		 exported outside their own package in golang

// PrintError is a utility function to handle errors gracefully.
func PrintError(errorMessage string, err error) {
	if err != nil {
		log.Println(errorMessage, color.Ize(color.Red, err)) // using log package to add timestamp with the error
		os.Exit(1)
	}
}

func GetKubeConfig() (*rest.Config, error) {
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
		return nil, fmt.Errorf("failed to create Kubernetes config: %v", err)
	}

	return config, nil

}
