package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

func init() {
  rootCmd.AddCommand(versionCmd)
}

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print the version number of TK8",
  Long:  `All software has versions. This is TK8's`,
  Run: func(cmd *cobra.Command, args []string) {
    fmt.Println("0.0.1")
  },
}