package k8s

import (
	"context"
	"fmt"
	"io"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// PodInfo contains information about a pod
type PodInfo struct {
	Name          string
	Namespace     string
	Phase         corev1.PodPhase
	Restarts      int32
	Ready         bool
	ContainerInfo []ContainerInfo
}

// ContainerInfo contains information about a container
type ContainerInfo struct {
	Name         string
	Image        string
	Ready        bool
	RestartCount int32
	State        string
}

// GetPods retrieves pods from the specified namespace
func (c *Client) GetPods(ctx context.Context, namespace string) ([]PodInfo, error) {
	if namespace == "" {
		namespace = c.namespace
	}
	
	pods, err := c.clientset.CoreV1().Pods(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list pods: %w", err)
	}
	
	var podInfos []PodInfo
	for _, pod := range pods.Items {
		podInfo := PodInfo{
			Name:      pod.Name,
			Namespace: pod.Namespace,
			Phase:     pod.Status.Phase,
		}
		
		// Calculate readiness and restarts
		ready := true
		totalRestarts := int32(0)
		for _, cs := range pod.Status.ContainerStatuses {
			if !cs.Ready {
				ready = false
			}
			totalRestarts += cs.RestartCount
			
			state := "running"
			if cs.State.Waiting != nil {
				state = cs.State.Waiting.Reason
			} else if cs.State.Terminated != nil {
				state = cs.State.Terminated.Reason
			}
			
			podInfo.ContainerInfo = append(podInfo.ContainerInfo, ContainerInfo{
				Name:         cs.Name,
				Image:        cs.Image,
				Ready:        cs.Ready,
				RestartCount: cs.RestartCount,
				State:        state,
			})
		}
		
		podInfo.Ready = ready
		podInfo.Restarts = totalRestarts
		
		podInfos = append(podInfos, podInfo)
	}
	
	return podInfos, nil
}

// GetPod retrieves a specific pod
func (c *Client) GetPod(ctx context.Context, name, namespace string) (*corev1.Pod, error) {
	if namespace == "" {
		namespace = c.namespace
	}
	
	pod, err := c.clientset.CoreV1().Pods(namespace).Get(ctx, name, metav1.GetOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to get pod %s: %w", name, err)
	}
	
	return pod, nil
}

// GetPodLogs retrieves logs from a pod
func (c *Client) GetPodLogs(ctx context.Context, podName, containerName, namespace string, tailLines int64) (string, error) {
	if namespace == "" {
		namespace = c.namespace
	}
	
	opts := &corev1.PodLogOptions{
		Container: containerName,
	}
	
	if tailLines > 0 {
		opts.TailLines = &tailLines
	}
	
	req := c.clientset.CoreV1().Pods(namespace).GetLogs(podName, opts)
	podLogs, err := req.Stream(ctx)
	if err != nil {
		return "", fmt.Errorf("failed to get logs: %w", err)
	}
	defer podLogs.Close()
	
	buf := new(io.Reader)
	*buf = podLogs
	
	bytes, err := io.ReadAll(podLogs)
	if err != nil {
		return "", fmt.Errorf("failed to read logs: %w", err)
	}
	
	return string(bytes), nil
}

// GetEvents retrieves events for a namespace
func (c *Client) GetEvents(ctx context.Context, namespace string) (*corev1.EventList, error) {
	if namespace == "" {
		namespace = c.namespace
	}
	
	events, err := c.clientset.CoreV1().Events(namespace).List(ctx, metav1.ListOptions{})
	if err != nil {
		return nil, fmt.Errorf("failed to list events: %w", err)
	}
	
	return events, nil
}

// DeletePod deletes a pod
func (c *Client) DeletePod(ctx context.Context, name, namespace string) error {
	if namespace == "" {
		namespace = c.namespace
	}
	
	err := c.clientset.CoreV1().Pods(namespace).Delete(ctx, name, metav1.DeleteOptions{})
	if err != nil {
		return fmt.Errorf("failed to delete pod %s: %w", name, err)
	}
	
	return nil
}
