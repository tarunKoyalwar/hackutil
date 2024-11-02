package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
	"github.com/tarunKoyalwar/hackutil/pkg/tools/burpconfig"
)

var rootCmd = &cobra.Command{
	Use:   "hackutil",
	Short: "A collection of utility tools",
	Long:  `hackutil is a collection of various utility tools for different purposes`,
}

func init() {
	rootCmd.AddCommand(burpconfig.NewBurpConfigCmd())
}

func main() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintf(os.Stderr, "Error: %v\n", err)
		os.Exit(1)
	}
}
