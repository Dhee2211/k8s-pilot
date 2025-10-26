package k8s

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

// Client wraps the Kubernetes client
type Client struct {
	clientset *kubernetes.Clientset
	config    *rest.Config
	namespace string
}

// NewClient creates a new Kubernetes client
func NewClient(namespace string) (*Client, error) {
	config, err := getConfig()
	if err != nil {
		return nil, fmt.Errorf("failed to get kubernetes config: %w", err)
	}
	
	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, fmt.Errorf("failed to create kubernetes client: %w", err)
	}
	
	if namespace == "" {
		namespace = "default"
	}
	
	return &Client{
		clientset: clientset,
		config:    config,
		namespace: namespace,
	}, nil
}

// getConfig returns the Kubernetes config from kubeconfig or in-cluster
func getConfig() (*rest.Config, error) {
	// Try in-cluster config first
	config, err := rest.InClusterConfig()
	if err == nil {
		return config, nil
	}
	
	// Fall back to kubeconfig
	var kubeconfig string
	if os.Getenv("KUBECONFIG") != "" {
		kubeconfig = os.Getenv("KUBECONFIG")
	} else if home := homedir.HomeDir(); home != "" {
		kubeconfig = filepath.Join(home, ".kube", "config")
	} else {
		return nil, fmt.Errorf("unable to find kubeconfig")
	}
	
	config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	if err != nil {
		return nil, fmt.Errorf("failed to build config from kubeconfig: %w", err)
	}
	
	return config, nil
}

// Clientset returns the underlying Kubernetes clientset
func (c *Client) Clientset() *kubernetes.Clientset {
	return c.clientset
}

// Config returns the Kubernetes config
func (c *Client) Config() *rest.Config {
	return c.config
}

// Namespace returns the current namespace
func (c *Client) Namespace() string {
	return c.namespace
}

// SetNamespace sets the namespace for operations
func (c *Client) SetNamespace(namespace string) {
	c.namespace = namespace
}

// DetectClusterType detects the type of Kubernetes cluster
func (c *Client) DetectClusterType(ctx context.Context) (ClusterType, error) {
	// Try to detect cluster type from API server
	version, err := c.clientset.Discovery().ServerVersion()
	if err != nil {
		return ClusterTypeUnknown, err
	}
	
	// Check for GKE
	if version.GitVersion != "" && len(version.GitVersion) > 0 {
		// GKE has specific version patterns
		// EKS has specific version patterns
		// AKS has specific version patterns
		// This is a simplified detection
	}
	
	// For now, return generic
	return ClusterTypeGeneric, nil
}

// ClusterType represents the type of Kubernetes cluster
type ClusterType string

const (
	ClusterTypeGKE     ClusterType = "gke"
	ClusterTypeEKS     ClusterType = "eks"
	ClusterTypeAKS     ClusterType = "aks"
	ClusterTypeK3s     ClusterType = "k3s"
	ClusterTypeKind    ClusterType = "kind"
	ClusterTypeGeneric ClusterType = "generic"
	ClusterTypeUnknown ClusterType = "unknown"
)
