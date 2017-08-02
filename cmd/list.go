// Copyright © 2017 The Kubicorn Authors
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

// Copyright © 2017 The Kubicorn Authors
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

	"github.com/kris-nova/kubicorn/logger"
	"github.com/kris-nova/kubicorn/state"
	"github.com/kris-nova/kubicorn/state/fs"
	"github.com/spf13/cobra"
)

var listCmd = &cobra.Command{
	Use:   "list",
	Short: "List available states",
	Long:  `List the states available in the _state directory`,
	Run: func(cmd *cobra.Command, args []string) {
		err := RunList(lo)
		if err != nil {
			logger.Critical(err.Error())
			os.Exit(1)
		}
	},
}

type ListOptions struct {
	Options
	Profile string
}

var lo = &ListOptions{}

func init() {
	listCmd.Flags().StringVarP(&lo.StateStore, "state-store", "s", strEnvDef("KUBICORN_STATE_STORE", "fs"), "The state store type to use for the cluster")
	listCmd.Flags().StringVarP(&lo.StateStorePath, "state-store-path", "S", strEnvDef("KUBICORN_STATE_STORE_PATH", "./_state"), "The state store path to use")
	RootCmd.AddCommand(listCmd)

}

func RunList(options *ListOptions) error {
	options.StateStorePath = expandPath(options.StateStorePath)

	var stateStore state.ClusterStorer
	switch options.StateStore {
	case "fs":
		logger.Info("Selected [fs] state store")
		stateStore = fs.NewFileSystemStore(&fs.FileSystemStoreOptions{
			BasePath: options.StateStorePath,
		})
	}

	clusters, err := stateStore.List()
	if err != nil {
		return fmt.Errorf("Unable to list clusters: %v", err)
	}
	for _, cluster := range clusters {
		logger.Always(cluster)
	}

	return nil
}
