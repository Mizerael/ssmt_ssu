package searchconfig

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"
)

type Config struct {
	Count        int
	Bosting      float64
	IndexRebuild bool
}

var conf Config

var rootCmd = &cobra.Command{
	Use:   "ssmt-ssu",
	Short: "show hide message",
	Long:  `Solution second part of first lab on course infoSec SSU`,

	Run: func(cmd *cobra.Command, args []string) {

	},
}

func init() {
	rootCmd.Flags().IntVarP(&conf.Count, "count", "c", 10,
		"search result count")
	rootCmd.Flags().Float64VarP(&conf.Bosting, "boost", "b", 2.0,
		"boosting power var")
	rootCmd.Flags().BoolVarP(&conf.IndexRebuild, "index-rebuild", "R", false,
		"rebuild index if set True")
}

func Execute() *Config {
	err := rootCmd.Execute()
	if err != nil {
		fmt.Printf("err: %v\n", err)
		os.Exit(1)
	}
	return &conf
}
