package kubeclient

import (
	"context"
	"fmt"
	"log"
	"time"

	"golang.org/x/exp/slices"

	utils "github.com/cloud-sky-ops/ice-kube/internal"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// DeleteResources connects to the Kubernetes API and cleans up completed pods, PVCs, and LoadBalancers
func DeleteResources(clusterName string, deleteBeforeHours int, namespace string, dryRun bool) (string, error) {

	clientset, err := utils.GetClientSet()

	if err != nil {
		return "", fmt.Errorf("failed to create Kubernetes client: %v", err)
	}

	pods, err := ScanResources(clusterName, namespace)

	if err != nil {
		return "", fmt.Errorf("failed to list pods: %v", err)
	}

	maxPermittedTime := time.Now().Add(time.Duration(deleteBeforeHours) * time.Second) // reduce 24 hours in current time to set maxPermittedTime
	deletedPods := 0

	podsToDelete := []string{}

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
				if dryRun {
					podsToDelete = append(podsToDelete, pod.Name)
					deletedPods++
				} else {
					podsToDelete = append(podsToDelete, pod.Name)
					err := clientset.CoreV1().Pods(namespace).Delete(context.TODO(), pod.Name, metav1.DeleteOptions{})
					if err == nil {
						log.Printf("Deleted completed pod: %s (namespace: %s, finished at: %s)", pod.Name, namespace, latestFinishTime)
						deletedPods++
					}
				}
			}
		}
	}

	// Get updated list of pods data across all namespaces in the cluster, after deletion.
	if !dryRun && deletedPods != 0 {
		pods, err = clientset.CoreV1().Pods(namespace).List(context.TODO(), metav1.ListOptions{})
		if err != nil {
			return "", fmt.Errorf("failed to list pods: %v", err)
		}

		fmt.Printf("Found %d pods in namespace %s after deletion\n", len(pods.Items), namespace)
	}

	// Delete unbouded PVCs
	pvcs, err := clientset.CoreV1().PersistentVolumeClaims(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list PVCs: %v", err)
	}

	deletedPVCs := 0
	PVCstoDelete := []string{}
	for _, pvc := range pvcs.Items {
		if pvc.Status.Phase == "Bound" {
			// Check if PVC is still linked to a pod
			isUsed := false
			for _, pod := range pods.Items {
				if !slices.Contains(podsToDelete, pod.Name) {
					for _, volume := range pod.Spec.Volumes {
						if volume.PersistentVolumeClaim != nil && volume.PersistentVolumeClaim.ClaimName == pvc.Name {
							isUsed = true
							break
						}
					}
				}
			}

			if !isUsed {
				if dryRun {
					PVCstoDelete = append(PVCstoDelete, pvc.Name)
					deletedPVCs++
				} else {
					PVCstoDelete = append(PVCstoDelete, pvc.Name)
					err := clientset.CoreV1().PersistentVolumeClaims(namespace).Delete(context.TODO(), pvc.Name, metav1.DeleteOptions{})
					if err == nil {
						log.Printf("Deleted unbound PVC: %s (namespace: %s)", pvc.Name, namespace)
						deletedPVCs++
					}
				}
			}
		}
	}

	// Delete unused loadBalancer Service for which the pods are deleted, these will have no endpoints or sub-endpoints
	services, err := clientset.CoreV1().Services(namespace).List(context.TODO(), metav1.ListOptions{})
	if err != nil {
		return "", fmt.Errorf("failed to list services: %v", err)
	}

	deletedServices := 0
	servicesToDelete := []string{}
	for _, svc := range services.Items {
		if svc.Spec.Type == "LoadBalancer" {
			// Check if any pod is linked to the service
			endpoints, _ := clientset.CoreV1().Endpoints(namespace).Get(context.TODO(), svc.Name, metav1.GetOptions{})
			if endpoints == nil {
				if dryRun {
					servicesToDelete = append(servicesToDelete, svc.Name)
					deletedServices++
				} else {
					servicesToDelete = append(servicesToDelete, svc.Name)
					err := clientset.CoreV1().Services(namespace).Delete(context.TODO(), svc.Name, metav1.DeleteOptions{})
					if err == nil {
						log.Printf("Deleted unused LoadBalancer service: %s (namespace: %s)", svc.Name, namespace)
						deletedServices++
					}
				}
			} else {
				if  len(endpoints.Subsets) == 0 {
					if dryRun {
						servicesToDelete = append(servicesToDelete, svc.Name)
						deletedServices++
					} else {
						servicesToDelete = append(servicesToDelete, svc.Name)
						err := clientset.CoreV1().Services(namespace).Delete(context.TODO(), svc.Name, metav1.DeleteOptions{})
						if err == nil {
							log.Printf("Deleted unused LoadBalancer service: %s (namespace: %s)", svc.Name, namespace)
							deletedServices++
						}
					}
				} else {
					for _, subset := range endpoints.Subsets {
						for _, NotReadyAddress := range subset.NotReadyAddresses {
							fmt.Printf("%s\n",NotReadyAddress.TargetRef.Name)
							if slices.Contains(podsToDelete, NotReadyAddress.TargetRef.Name) {
								if dryRun {
									servicesToDelete = append(servicesToDelete, svc.Name)
									deletedServices++
								} else {
									servicesToDelete = append(servicesToDelete, svc.Name)
									err := clientset.CoreV1().Services(namespace).Delete(context.TODO(), svc.Name, metav1.DeleteOptions{})
									if err == nil {
										log.Printf("Deleted unused LoadBalancer service: %s (namespace: %s)", svc.Name, namespace)
										deletedServices++
									}
								}
							}
						}
					}
				}
			}	
		}
	}
	
	if dryRun {
		return fmt.Sprintf("%d Pods to delete: %s\n %d PVCs to delete: %s\n %d LoadBalancer services to delete: %s\n", deletedPods, podsToDelete, deletedPVCs, PVCstoDelete, deletedServices, servicesToDelete), nil
	} else {
		return fmt.Sprintf("Deleted %d Pods: %s\n Deleted %d PVCs: %s\n Deleted %d LoadBalancer services: %s", deletedPods, podsToDelete, deletedPVCs, PVCstoDelete, deletedServices, servicesToDelete), nil
	}
}
