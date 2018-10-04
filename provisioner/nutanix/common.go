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

type Nutanix struct {
}

func (p Nutanix) Init(args []string) {
	cluster.NotImplemented()
}

func (p Nutanix) Setup(args []string) {
	cluster.NotImplemented()
}

func (p Nutanix) Scale(args []string) {
	cluster.NotImplemented()

}

func (p Nutanix) Reset(args []string) {
	cluster.NotImplemented()

}

func (p Nutanix) Remove(args []string) {
	cluster.NotImplemented()

}

func (p Nutanix) Upgrade(args []string) {
	cluster.NotImplemented()
}

func (p Nutanix) Destroy(args []string) {
	cluster.NotImplemented()
}

func NewNutanix() cluster.Provisioner {
	provisioner := new(Nutanix)
	return provisioner
}
