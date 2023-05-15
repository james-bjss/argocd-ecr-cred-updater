package k8s

import (
	"context"
	"encoding/json"
	"fmt"
	v1 "k8s.io/api/core/v1"
	"k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/types"
	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
)

type K8Client struct {
	client    kubernetes.Interface
	namespace string
}

func New(client kubernetes.Interface, namespace string) K8Client {
	return K8Client{
		client:    client,
		namespace: namespace,
	}
}

func ConfigureClientWithConfig(kubeConfig *string) (*kubernetes.Clientset, error) {
	// use the current context in kubeconfig
	config, err := clientcmd.BuildConfigFromFlags("", *kubeConfig)
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func ConfigureInClusterClient() (*kubernetes.Clientset, error) {
	config, err := rest.InClusterConfig()
	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func (k K8Client) PatchSecret(secretName string, secret string) error {
	_, err := k.client.CoreV1().Secrets(k.namespace).Get(context.TODO(), secretName, metav1.GetOptions{})
	if errors.IsNotFound(err) {
		return fmt.Errorf("secret %s in namespace %s not found\n", secret, k.namespace)
	} else if err != nil {
		return err
	}

	// patch secret
	patch := v1.Secret{
		Data: map[string][]byte{
			"password": []byte(secret),
		},
	}
	payloadBytes, err := json.Marshal(patch)
	if err != nil {
		return err
	}

	_, err = k.client.CoreV1().Secrets(k.namespace).Patch(context.TODO(), secretName, types.MergePatchType, payloadBytes, metav1.PatchOptions{})
	return err
}
