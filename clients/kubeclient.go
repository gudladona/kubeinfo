package clients

import (
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
)

//NewClientSet creates a Clientset with Incluster Configuration which
// uses uses the Service Account token mounted inside the Pod
func NewClientSet() (*kubernetes.Clientset, error) {
	// creates the in-cluster config
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	// creates the clientset
	return kubernetes.NewForConfig(config)
}
