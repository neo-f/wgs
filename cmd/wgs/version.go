package main

import (
	"fmt"

	"wgs"

	"github.com/spf13/cobra"
)

var cmdVersion = &cobra.Command{
	Use:   "version",
	Short: "print the version info",
	Run: func(*cobra.Command, []string) {
		fmt.Println("version:", wgs.Version)
		fmt.Println("build time:", wgs.BuildTime)
		fmt.Println("build go version:", wgs.GoVersion)
	},
}

func init() { // nolint
	rootCmd.AddCommand(cmdVersion)
}
