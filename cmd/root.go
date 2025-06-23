/*
Copyright © 2025 Andrii Ivanov <bergshrund@gmail.com>
*/
package cmd

import (
	"os"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/log"
	"github.com/spf13/cobra"
)

var logLevel string

var rootCmd = &cobra.Command{
	Use:   "k8s-controller",
	Short: "A template for Kubernetes controller",
	Long:  "A template for Kubernetes controller",
	Run: func(cmd *cobra.Command, args []string) {
		configureLogLevel(getLogLevel(logLevel))
		log.Info().Msgf("k8s-controller %s started", appVersion)
	},
}

func getLogLevel(logLevel string) zerolog.Level {
	switch logLevel {
	case "trace":
		return zerolog.TraceLevel
	case "debug":
		return zerolog.DebugLevel
	case "info":
		return zerolog.InfoLevel
	case "warn":
		return zerolog.WarnLevel
	case "error":
		return zerolog.ErrorLevel
	case "fatal":
		return zerolog.FatalLevel
	default:
		log.Warn().Msgf("Unknown log level %s, defaulting to info", logLevel)
		return zerolog.InfoLevel
	}
}

func configureLogLevel(logLevel zerolog.Level) {
	zerolog.TimeFieldFormat = zerolog.TimeFormatUnix
	zerolog.SetGlobalLevel(logLevel)
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(1)
	}
}

func init() {
	rootCmd.PersistentFlags().StringVar(&logLevel, "log-level", "info", "Log level (trace, debug, info, warn, error, fatal)")
}
