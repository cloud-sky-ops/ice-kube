package kubeclient

import (
	"context"
	"errors"
	"fmt"

	utils "github.com/cloud-sky-ops/ice-kube/internal"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

func ScanResources(clusterName string, namespace string) (*v1.PodList, error) {
	clientset, err := utils.GetClientSet()

	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	namespaceList, err := clientset.CoreV1().Namespaces().List(context.TODO(), metav1.ListOptions{})

	if err != nil {
		fmt.Print(err)
	}
	namespaceFound := false

	if namespace == "" {
		namespaceFound = true
	}

	if !namespaceFound {
		for _, ns := range namespaceList.Items {
			if ns.Name == namespace {
				namespaceFound = true
			}
		}
	}

	if !namespaceFound {
		message := "NAMESPACE NOT FOUNT IN CLUSTER"
		namespaceError := errors.New("RE-CHECK INPUT --namespace")
		utils.PrintError(message, namespaceError)
	}

	// Get pods data across all namespaces in the cluster
	pods, err := clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %v", err)
	}
	return pods, nil
}
