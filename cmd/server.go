/*
Copyright © 2025 Andrii Ivanov <bergshrund@gmail.com>
*/
package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"

	"k8s-controller-tmpl/pkg/informer"
)

var (
	serverPort          int
	inClusterConfigFlag bool = false
)

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start http server",
	Run: func(cmd *cobra.Command, args []string) {
		configureLogLevel(getLogLevel(logLevel))

		ctx := context.Background()

		clientset, err := getKubeClient(kubeconfig)
		if err != nil {
			log.Error().Err(err).Msg("Error creating Kubernetes client")
			os.Exit(1)
		}

		go informer.StartInformer(ctx, clientset, namespace)

		addr := fmt.Sprintf(":%d", serverPort)
		log.Info().Msgf("Starting server on port %s", addr)

		router := gin.Default()

		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{
				"error":  "Endpoint not implemented",
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
			})
		})

		if err := http.ListenAndServe(addr, router); err != nil {
			log.Error().Err(err).Msg("Error starting http server")
			os.Exit(1)
		}

		http.ListenAndServe(addr, router)
	},
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVar(&serverPort, "port", 8080, "Port to run the server on")
	serverCmd.Flags().StringVar(&kubeconfig, "kubeconfig", getKubeconfigPath(), "Path to the kubeconfig file")
	serverCmd.Flags().StringVar(&namespace, "namespace", "default", "Kubernetes namespace")
	serverCmd.Flags().BoolVar(&inClusterConfigFlag, "in-cluster", false, "Use in-cluster Kubernetes config")
}
