package cmd

import (
	"fmt"
	"os"

	"github.com/ezhigval/algo-sandbox/internal/algo/graph"
	"github.com/ezhigval/algo-sandbox/internal/algo/lru"
	"github.com/spf13/cobra"
)

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
}

var rootCmd = &cobra.Command{
	Use:   "algo",
	Short: "Run algorithm demos from the terminal",
}

var benchLRUCmd = &cobra.Command{
	Use:   "bench-lru",
	Short: "Quick LRU smoke demo",
	Run: func(cmd *cobra.Command, args []string) {
		c := lru.New(2)
		c.Put("x", "1")
		c.Put("y", "2")
		c.Get("x")
		c.Put("z", "3")
		if _, ok := c.Get("y"); ok {
			fmt.Println("FAIL: y should be evicted")
			os.Exit(1)
		}
		fmt.Println("LRU ok")
	},
}

var graphDemoCmd = &cobra.Command{
	Use:   "graph-demo",
	Short: "BFS shortest path demo",
	Run: func(cmd *cobra.Command, args []string) {
		g := graph.Graph{
			"A": {"B", "C"},
			"B": {"D"},
			"C": {"D"},
			"D": {},
		}
		path, ok := g.BFS("A", "D")
		if !ok {
			fmt.Println("no path")
			os.Exit(1)
		}
		fmt.Println("path:", path)
	},
}

func init() {
	rootCmd.AddCommand(benchLRUCmd)
	rootCmd.AddCommand(graphDemoCmd)
}
