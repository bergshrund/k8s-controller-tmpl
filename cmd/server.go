/*
Copyright © 2025 Andrii Ivanov <bergshrund@gmail.com>
*/
package cmd

import (
	"context"
	"fmt"
	"net/http"
	"os"

	"k8s-controller-tmpl/pkg/controller"
	"k8s-controller-tmpl/pkg/informer"

	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
	runtime "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/cache"
	"sigs.k8s.io/controller-runtime/pkg/metrics/server"
)

var (
	serverPort          int
	inClusterConfigFlag bool = false
	metricsPort         int
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

		managedNamespace := map[string]cache.Config{}
		managedNamespace[namespace] = cache.Config{}

		controllerManager, err := runtime.NewManager(runtime.GetConfigOrDie(), runtime.Options{
			Metrics: server.Options{BindAddress: fmt.Sprintf(":%d", metricsPort)},
			Cache: cache.Options{
				DefaultNamespaces: managedNamespace,
			},
		})

		if err != nil {
			log.Error().Err(err).Msg("Error creating controller-runtime manager")
			os.Exit(1)
		}

		if err := controller.AddDeploymentController(controllerManager); err != nil {
			log.Error().Err(err).Msg("Failed to add deployment controller")
			os.Exit(1)
		}

		go func() {
			log.Info().Msg("Starting controller-runtime manager...")
			if err := controllerManager.Start(cmd.Context()); err != nil {
				log.Error().Err(err).Msg("Error starting controller-runtime manager")
				os.Exit(1)
			}
		}()

		addr := fmt.Sprintf(":%d", serverPort)
		log.Info().Msgf("Starting server on port %s", addr)

		router := gin.Default()
		router.Use(gin.Recovery())
		router.Use(requestID())

		router.NoRoute(func(c *gin.Context) {
			c.JSON(http.StatusNotImplemented, gin.H{
				"error":  "Endpoint not implemented",
				"path":   c.Request.URL.Path,
				"method": c.Request.Method,
			})
		})

		router.GET("/deployments", func(c *gin.Context) {
			deployments := informer.GetDeploymentNames()

			requestID, _ := c.Get("X-Request-ID")

			logger := log.With().Str("request_id", requestID.(string)).Logger()
			logger.Info().Msgf("Deployments: %v", deployments)

			c.JSON(http.StatusOK, gin.H{
				"status":      "ok",
				"deployments": deployments,
			})
		})

		if err := http.ListenAndServe(addr, router); err != nil {
			log.Error().Err(err).Msg("Error starting http server")
			os.Exit(1)
		}

		http.ListenAndServe(addr, router)
	},
}

func requestID() gin.HandlerFunc {
	return func(c *gin.Context) {
		requestID := c.GetHeader("X-Request-ID")

		if requestID == "" {
			requestID = uuid.New().String()
		}

		c.Set("X-Request-ID", requestID)
		c.Writer.Header().Set("X-Request-ID", requestID)
		c.Next()
	}
}

func init() {
	rootCmd.AddCommand(serverCmd)
	serverCmd.Flags().IntVar(&serverPort, "port", 8080, "Port to run the server on")
	serverCmd.Flags().StringVar(&kubeconfig, "kubeconfig", getKubeconfigPath(), "Path to the kubeconfig file")
	serverCmd.Flags().StringVar(&namespace, "namespace", "default", "Kubernetes namespace")
	serverCmd.Flags().BoolVar(&inClusterConfigFlag, "in-cluster", false, "Use in-cluster Kubernetes config")
	serverCmd.Flags().IntVar(&metricsPort, "metrics-port", 8081, "Port for controller manager metrics")
}
