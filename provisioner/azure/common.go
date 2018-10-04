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

package provisioner

import "github.com/kubernauts/tk8/internal/cluster"

type Azure struct {
}

func (p Azure) Init(args []string) {
	cluster.NotImplemented()
}

func (p Azure) Setup(args []string) {
	cluster.NotImplemented()
}

func (p Azure) Scale(args []string) {
	cluster.NotImplemented()

}

func (p Azure) Reset(args []string) {
	cluster.NotImplemented()

}

func (p Azure) Remove(args []string) {
	cluster.NotImplemented()

}

func (p Azure) Upgrade(args []string) {
	cluster.NotImplemented()
}

func (p Azure) Destroy(args []string) {
	cluster.NotImplemented()
}

func NewAzure() cluster.Provisioner {
	provisioner := new(Azure)
	return provisioner
}
