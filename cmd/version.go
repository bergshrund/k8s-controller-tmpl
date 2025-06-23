/*
Copyright © 2025 Andrii Ivanov <bergshrund@gmail.com>
*/
package cmd

import (
	"fmt"

	"github.com/spf13/cobra"
)

var appVersion = "Version"

// versionCmd represents the version command
var versionCmd = &cobra.Command{
	Use:   "version",
	Short: "Current version of the application",
	Long:  "Current version of the application",
	Run: func(cmd *cobra.Command, args []string) {
		fmt.Println(appVersion)
	},
}

func init() {
	rootCmd.AddCommand(versionCmd)
}
