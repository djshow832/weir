// Copyright 2020 Ipalfish, Inc.
// Copyright 2022 PingCAP, Inc.
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

package configcenter

import (
	"github.com/djshow832/weir/pkg/config"
	"github.com/pingcap/errors"
)

const (
	ConfigCenterTypeFile = "file"
	ConfigCenterTypeEtcd = "etcd"
)

type ConfigCenter interface {
	GetNamespace(ns string) (*config.Namespace, error)
	ListAllNamespace() ([]*config.Namespace, error)
}

func CreateConfigCenter(cfg config.ConfigCenter) (ConfigCenter, error) {
	switch cfg.Type {
	case ConfigCenterTypeFile:
		return CreateFileConfigCenter(cfg.ConfigFile.Path)
	case ConfigCenterTypeEtcd:
		return CreateEtcdConfigCenter(cfg.ConfigEtcd)
	default:
		return nil, errors.New("invalid config center type")
	}
}
