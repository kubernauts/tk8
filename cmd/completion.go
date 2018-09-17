// Copyright © 2018 The TK8 Authors.
//
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
//
//     http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.

package cmd

import (
	"fmt"
	"os"

	"github.com/kubernauts/tk8/internal/cluster"
	"github.com/spf13/cobra"
)

// completionCmd represents the completion command.
var completionCmd = &cobra.Command{
	Use:   "completion",
	Short: "Generates shell autocompletion script for either bash or zsh.",
	Long: `A longer description that spans multiple lines and likely contains examples
and usage of using your command. For example:

Cobra is a CLI library for Go that empowers applications.
This application is a tool to generate the needed files
to quickly create a Cobra application.`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) == 0 {
			cmd.Help()
			os.Exit(0)
		}
	},
}

// bashCompletion represents the bash sub-command of the completion command.
var bashCompletion = &cobra.Command{
	Use:   "bash",
	Short: "Generates shell autocompletion script for bash.",
	Long:  `It will produce the bash completion script which can later be used for the autocompletion of commands in Bash.`,
	Run: func(cmd *cobra.Command, args []string) {
		script, err := os.OpenFile("tk8.sh", os.O_CREATE|os.O_WRONLY, 0600)
		cluster.ErrorCheck("Error creating autocompletion script file.", err)
		err = rootCmd.GenBashCompletion(script)
		fmt.Printf("Successfully created the Bash completion script. Move the 'tk8.sh' file under /etc/bash_completion.d/ folder and login again.")
	},
}

// zshCompletion represents the zsh sub-command of the completion command.
var zshCompletion = &cobra.Command{
	Use:   "zsh",
	Short: "Generates shell autocompletion script for zsh.",
	Long:  `It will produce the bash completion script which can later be used for the autocompletion of commands in Zsh.`,
	Run: func(cmd *cobra.Command, args []string) {
		script, err := os.OpenFile("tk8.plugin.zsh", os.O_CREATE|os.O_WRONLY, 0600)
		cluster.ErrorCheck("Error creating autocompletion script file.", err)

		fmt.Fprintf(script, "__tk8_tool_complete() {\n")
		err = rootCmd.GenZshCompletion(script)
		cluster.ErrorCheck("Zsh issue", err)
		fmt.Fprintf(script, "}\ncompdef __tk8_tool_complete tk8")
		fmt.Printf("Successfully created the Zsh plugin. Move the 'tk8.plugin.zsh' file under your plugins folder and login again.")
	},
}

func init() {
	// Add the parent completion command.
	rootCmd.AddCommand(completionCmd)
	// Add the child bash and zsh sub-commands.
	completionCmd.AddCommand(bashCompletion)
	completionCmd.AddCommand(zshCompletion)
}
