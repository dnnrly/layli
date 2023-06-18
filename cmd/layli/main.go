package main

import (
	"fmt"
	"os"

	"github.com/dnnrly/layli"
	"github.com/spf13/cobra"
)

func main() {
	err := Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	var output string
	var rootCmd = &cobra.Command{
		Use:   "layli [flags] [layout file]",
		Short: "",
		Long:  ``,
		Args:  cobra.MatchAll(cobra.ExactArgs(1), cobra.OnlyValidArgs),
		RunE: func(cmd *cobra.Command, args []string) error {
			f, err := os.Open(args[0])
			if err != nil {
				return fmt.Errorf("opening input: %w", err)
			}

			d, err := layli.NewDiagram(
				f,
				func(output string) error {
					name := fmt.Sprintf("%s.svg", args[0])
					return os.WriteFile(name, []byte(output), 0644)
				},
			)
			if err != nil {
				return fmt.Errorf("creating diagram: %w", err)
			}

			err = d.Draw()
			if err != nil {
				return fmt.Errorf("drawing diagram: %w", err)
			}

			return nil
		},
	}

	rootCmd.PersistentFlags().StringVarP(&output, "output", "o", "", "output file or directory/")

	return rootCmd.Execute()
}
