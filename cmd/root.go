package cmd

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

var (
	inputFile string
	output    string
)

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() {
	var rootCmd = &cobra.Command{
		Use:   "layli [flags] [layout file]",
		Short: "",
		Long:  ``,
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			_, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("opening input: %w", err)
			}

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output file or directory/")

	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}
