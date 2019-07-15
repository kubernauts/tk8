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
	"log"
	"os"

	"github.com/kubernauts/tk8/api/server"
	"github.com/kubernauts/tk8/pkg/common"
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
		if err := server.StartTK8API("tk8", common.REST_API_PORT); err != nil {
			logrus.Printf("Unable to start cluster API server: %v", err)
		}

		// daemon does not exit
		select {}
	},
	PreRun: func(cmd *cobra.Command, args []string) {

		if common.REST_API_PORT <= 0 {
			log.Fatal("Port number cannot be zero")
		}
		switch common.REST_API_STORAGE {
		case "local":
			isExists := checkStoragePath(common.REST_API_STORAGEPATH)
			if !isExists {
				log.Fatalf("Storage path [ %s ] either doesnt exist or there is an error", common.REST_API_STORAGEPATH)
			}
		case "s3":
		default:
			log.Fatal("storage flag accepts local or s3 as valid values")
		}
	},
}

func init() {
	daemonCmd.Flags().Uint16VarP(&common.REST_API_PORT, "port", "p", 8091, "Port number for the Tk8 rest api")
	daemonCmd.Flags().StringVarP(&common.REST_API_STORAGE, "config-store", "s", "local", "Storage for config files - local or s3")
	daemonCmd.Flags().StringVarP(&common.REST_API_STORAGEPATH, "confif-store path", "a", ".", "Storage location for config files - directory path for local , bucket name for s3")

	rootCmd.AddCommand(daemonCmd)
}

func checkStoragePath(path string) bool {
	exists := true
	src, err := os.Stat(path)
	if os.IsNotExist(err) {
		exists = false
	}
	if os.IsExist(err) {
		if src.Mode().IsDir() {
			exists = true
		}
	}
	return exists
}
