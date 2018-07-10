package cmd

import (
  "fmt"

  "github.com/spf13/cobra"
)

var (
      VERSION = "0.0.3"
      GITCOMMIT = "HEAD"
      client := github.NewClient(nil)

)

var versionCmd = &cobra.Command{
  Use:   "version",
  Short: "Print the version of TK8",
  Long:  `All software has versions. This is TK8's`,
  Run: func(cmd *cobra.Command, args []string) {
  
    fmt.Println(VERSION + " (" + GITCOMMIT + ")")
  },
}

func init() {
	rootCmd.AddCommand(versionCmd)
}

