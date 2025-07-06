/*
Copyright © 2025 Andrii Ivanov <bergshrund@gmail.com>
*/
package cmd

import (
	"fmt"
	"net/http"
	"os"

	"github.com/gin-gonic/gin"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var serverPort int

// serverCmd represents the server command
var serverCmd = &cobra.Command{
	Use:   "server",
	Short: "Start http server",
	Run: func(cmd *cobra.Command, args []string) {
		configureLogLevel(getLogLevel(logLevel))

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
	serverCmd.PersistentFlags().IntVar(&serverPort, "port", 8080, "Port to run the server on")
}
