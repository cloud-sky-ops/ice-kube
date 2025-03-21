package kubeclient

import (
	"context"
	"fmt"
	"log"
	"time"
	utils "github.com/cloud-sky-ops/ice-kube/internal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/client-go/kubernetes"
)

// ScanCluster connects to the Kubernetes API and cleans up completed pods, PVCs, and LoadBalancers
func ScanCluster(clusterName string) (string, error) {

	config, err := utils.GetKubeConfig()

	if err != nil {
		return "", fmt.Errorf("failed to create Kubernetes config: %v", err)
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return "", fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	// Get pods data across all namespaces in the cluster
	pods, err := clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list pods: %v", err)
	}
	fmt.Printf("Found %d pods in cluster %s\n", len(pods.Items), clusterName)

	maxPermittedTime := time.Now().Add(-24 * time.Hour) // reduce 24 hours in current time to set maxPermittedTime
	deletedPods := 0

	for _, pod := range pods.Items {
		// Check if pod is exited
		if pod.Status.Phase == "Succeeded" {
			var latestFinishTime time.Time

			// Iterate over all container statuses and find the most recently finishedAt timestamp
			for _, cs := range pod.Status.ContainerStatuses {
				if cs.State.Terminated != nil {
					finishTime := cs.State.Terminated.FinishedAt.Time
					if finishTime.After(latestFinishTime) {
						latestFinishTime = finishTime
					}
				}
			}

			// Compare latestFinishTime with expiration threshold and delete resources fitting the set criteria
			// Anything before 24 hours from now should be deleted (Before maxPermittedTime)
			if !latestFinishTime.IsZero() && latestFinishTime.Before(maxPermittedTime) {
				err := clientset.CoreV1().Pods(pod.Namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
				if err == nil {
					log.Printf("Deleted completed pod: %s (namespace: %s, finished at: %s)", pod.Name, pod.Namespace, latestFinishTime)
					deletedPods++
				}
			}
		}
	}

	// Get updated list of pods data across all namespaces in the cluster, after deletion.
	pods, err = clientset.CoreV1().Pods("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list pods: %v", err)
	}

	fmt.Printf("Found %d pods in cluster %s after deletion\n", len(pods.Items), clusterName)

	// Delete unbouded PVCs
	pvcs, err := clientset.CoreV1().PersistentVolumeClaims("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list PVCs: %v", err)
	}

	deletedPVCs := 0
	for _, pvc := range pvcs.Items {
		if pvc.Status.Phase == "Bound" {
			// Check if PVC is still linked to a pod
			isUsed := false
			for _, pod := range pods.Items {
				for _, volume := range pod.Spec.Volumes {
					if volume.PersistentVolumeClaim != nil && volume.PersistentVolumeClaim.ClaimName == pvc.Name {
						isUsed = true
						break
					}
				}
			}

			if !isUsed {
				err := clientset.CoreV1().PersistentVolumeClaims(pvc.Namespace).Delete(context.TODO(), pvc.Name, metav1.DeleteOptions{})
				if err == nil {
					log.Printf("Deleted unbound PVC: %s (namespace: %s)", pvc.Name, pvc.Namespace)
					deletedPVCs++
				}
			}
		}
	}

	// Delete unused loadBalancer Service for which the pods are deleted, these will have no endpoints or sub-endpoints
	services, err := clientset.CoreV1().Services("").List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list services: %v", err)
	}

	deletedServices := 0
	for _, svc := range services.Items {
		if svc.Spec.Type == "LoadBalancer" {
			// Check if any pod is linked to the service
			endpoints, _ := clientset.CoreV1().Endpoints(svc.Namespace).Get(context.TODO(), svc.Name, metav1.GetOptions{})
			if endpoints == nil || len(endpoints.Subsets) == 0 {
				err := clientset.CoreV1().Services(svc.Namespace).Delete(context.TODO(), svc.Name, metav1.DeleteOptions{})
				if err == nil {
					log.Printf("Deleted unused LoadBalancer service: %s (namespace: %s)", svc.Name, svc.Namespace)
					deletedServices++
				}
			}
		}
	}

	return fmt.Sprintf("Deleted: %d pods, %d PVCs, %d LoadBalancer services", deletedPods, deletedPVCs, deletedServices), nil
}
