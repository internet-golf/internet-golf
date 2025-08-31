package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var rootCmd = &cobra.Command{
	Use:   "golf",
	Short: "A server",
}

func main() {
	rootCmd.AddCommand(&cobra.Command{
		Use:   "create-deployment",
		Short: "Creates a deployment",
		Run: func(cmd *cobra.Command, args []string) {
			fmt.Println("hello i am creating a deployment")
		},
	})

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}
