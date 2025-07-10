/*
Copyright © 2025 Andrii Ivanov <bergshrund@gmail.com>
*/
package cmd

import (
	"context"
	"fmt"
	"os"
	"path/filepath"

	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	"k8s.io/client-go/kubernetes"
	"k8s.io/client-go/rest"
	"k8s.io/client-go/tools/clientcmd"
	"k8s.io/client-go/util/homedir"
)

var (
	kubeconfig string
	namespace  string
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List Kubernetes resources in the specified namespace",
	Run: func(cmd *cobra.Command, args []string) {
		clientset, err := getKubeClient(kubeconfig)
		if err != nil {
			log.Error().Err(err).Msg("Error creating Kubernetes client")
			os.Exit(1)
		}

		deployments, err := clientset.AppsV1().Deployments(namespace).List(context.Background(), metav1.ListOptions{})

		if err != nil {
			log.Error().Err(err).Msg("Failed to list deployments")
			os.Exit(1)
		}
		fmt.Printf("Found %d deployments in '%s' namespace:\n", len(deployments.Items), namespace)
		for _, d := range deployments.Items {
			fmt.Println("-", d.Name)
		}

	},
}

func getKubeconfigPath() string {
	if kubeconfig := os.Getenv("KUBECONFIG"); kubeconfig != "" {
		return kubeconfig
	} else if home := homedir.HomeDir(); home != "" {
		return filepath.Join(home, ".kube", "config")
	} else {
		return ""
	}
}

func getKubeClient(kubeconfig string) (*kubernetes.Clientset, error) {

	var err error
	var config *rest.Config

	if inClusterConfigFlag {
		config, err = rest.InClusterConfig()
	} else {
		config, err = clientcmd.BuildConfigFromFlags("", kubeconfig)
	}

	if err != nil {
		return nil, err
	}

	clientset, err := kubernetes.NewForConfig(config)
	if err != nil {
		return nil, err
	}

	return clientset, nil
}

func init() {
	rootCmd.AddCommand(listCmd)
	listCmd.Flags().StringVar(&kubeconfig, "kubeconfig", getKubeconfigPath(), "Path to the kubeconfig file")
	listCmd.Flags().StringVar(&namespace, "namespace", "default", "Kubernetes namespace")
}
