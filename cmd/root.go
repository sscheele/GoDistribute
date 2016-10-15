package cmd

import "github.com/spf13/cobra"

//RootCmd is the root command for cobra
var RootCmd = &cobra.Command{
	Use:   "distribute",
	Short: "Distribute is a program to share files over a LAN",
	Run:   start,
	PreRunE: func(cmd *cobra.Command, args []string) error {
		if chunkSize != 1<<26 {
			chunkSize *= 1 << 20
		}
		return nil
	},
}

var (
	chunkSize   int64
	numBranches int
)

func init() {
	RootCmd.Flags().Int64VarP(&chunkSize, "chunk-size", "c", 1<<26, "File chunk size (MiB)")
	RootCmd.Flags().IntVarP(&numBranches, "num-branches", "n", 5, "Number of outgoing connections to make at once")
}

func start(cmd *cobra.Command, args []string) {

}
