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
	"os"

	"github.com/kubernauts/tk8/api/server"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var daemonCmd = &cobra.Command{
	Use:   "daemon-mode",
	Short: "Start tk8 in daemon mode",
	Long:  `Daemon mode for TK8 , listens on port ::8091`,
	Run: func(cmd *cobra.Command, args []string) {
		if len(args) > 0 {
			cmd.Help()
			os.Exit(0)
		}

		// Start the REST server as well as cli
		if err := server.StartTK8API("dummy", 8091); err != nil {
			logrus.Printf("Unable to start cluster API server: %v", err)
		}

		// daemon does not exit
		select {}
	},
}

func init() {
	rootCmd.AddCommand(daemonCmd)
}
