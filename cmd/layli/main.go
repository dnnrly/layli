package main

import (
	"fmt"
	"os"
	"strings"

	"github.com/dnnrly/layli"
	"github.com/dnnrly/layli/pathfinder/dijkstra"
	"github.com/spf13/cobra"
)

func main() {
	err := Execute()
	if err != nil {
		fmt.Fprintln(os.Stderr, err.Error())
		os.Exit(1)
	}
}

var newPathFinder = func(start, end dijkstra.Point) layli.PathFinder {
	fmt.Printf("Creating new path finder %s,%s\n", start, end)
	return dijkstra.NewPathFinder(start, end)
}

// Execute adds all child commands to the root command and sets flags appropriately.
// This is called by main.main(). It only needs to happen once to the rootCmd.
func Execute() error {
	var output string
	var showGrid bool

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

			if output == "" {
				output = strings.ReplaceAll(args[0], ".layli", "")
				output = fmt.Sprintf("%s.svg", output)
			}

			d, err := layli.NewDiagramFromFile(
				newPathFinder, f,
				func(data string) error {
					return os.WriteFile(output, []byte(data), 0644)
				},
				showGrid,
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
	rootCmd.PersistentFlags().BoolVar(&showGrid, "show-grid", false, "show the path grid dots (great for debugging)")

	return rootCmd.Execute()
}
