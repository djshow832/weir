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

package config

import "github.com/goccy/go-yaml"

const (
	DefaultClusterName = "default"

	MIN_SESSION_TIMEOUT = 600
)

type Proxy struct {
	Version      string       `yaml:"version"`
	Cluster      string       `yaml:"cluster"`
	ProxyServer  ProxyServer  `yaml:"proxy_server"`
	AdminServer  AdminServer  `yaml:"admin_server"`
	Log          Log          `yaml:"log"`
	Registry     Registry     `yaml:"registry"`
	ConfigCenter ConfigCenter `yaml:"config_center"`
	Security     Security     `yaml:"security"`
}

type ProxyServer struct {
	Addr           string `yaml:"addr"`
	MaxConnections uint32 `yaml:"max_connections"`
	SessionTimeout int    `yaml:"session_timeout"`
	StoragePath    string `yaml:"storage_path"`
	TCPKeepAlive   bool   `yaml:"tcp_keep_alive"`
	PDAddr         string `yaml:"pd_addrs"`
}

type AdminServer struct {
	Addr            string `yaml:"addr"`
	EnableBasicAuth bool   `yaml:"enable_basic_auth"`
	User            string `yaml:"user"`
	Password        string `yaml:"password"`
}

type Log struct {
	Level   string  `yaml:"level"`
	Format  string  `yaml:"format"`
	LogFile LogFile `yaml:"log_file"`
}

type LogFile struct {
	Filename   string `yaml:"filename"`
	MaxSize    int    `yaml:"max_size"`
	MaxDays    int    `yaml:"max_days"`
	MaxBackups int    `yaml:"max_backups"`
}

type Registry struct {
	Enable bool     `yaml:"enable"`
	Type   string   `yaml:"type"`
	Addrs  []string `yaml:"addrs"`
}

type ConfigCenter struct {
	Type       string     `yaml:"type"`
	ConfigFile ConfigFile `yaml:"config_file"`
	ConfigEtcd ConfigEtcd `yaml:"config_etcd"`
}

type ConfigFile struct {
	Path string `yaml:"path"`
}

type ConfigEtcd struct {
	Addrs    []string `yaml:"addrs"`
	BasePath string   `yaml:"base_path"`
	Username string   `yaml:"username"`
	Password string   `yaml:"password"`
	// If strictParse is disabled, parsing namespace error will be ignored when list all namespaces.
	StrictParse bool `yaml:"strict_parse"`
}

type Security struct {
	SSLCA           string   `toml:"ssl-ca" json:"ssl-ca"`
	SSLCert         string   `toml:"ssl-cert" json:"ssl-cert"`
	SSLKey          string   `toml:"ssl-key" json:"ssl-key"`
	ClusterSSLCA    string   `toml:"cluster-ssl-ca" json:"cluster-ssl-ca"`
	ClusterSSLCert  string   `toml:"cluster-ssl-cert" json:"cluster-ssl-cert"`
	ClusterSSLKey   string   `toml:"cluster-ssl-key" json:"cluster-ssl-key"`
	ClusterVerifyCN []string `toml:"cluster-verify-cn" json:"cluster-verify-cn"`
	MinTLSVersion   string   `toml:"tls-version" json:"tls-version"`
	RSAKeySize      int      `toml:"rsa-key-size" json:"rsa-key-size"`
}

func NewProxyConfig(data []byte) (*Proxy, error) {
	var cfg Proxy
	if err := yaml.Unmarshal(data, &cfg); err != nil {
		return nil, err
	}
	if err := cfg.Check(); err != nil {
		return nil, err
	}
	return &cfg, nil
}

func (cfg *Proxy) Check() error {
	if cfg.ProxyServer.SessionTimeout <= MIN_SESSION_TIMEOUT {
		cfg.ProxyServer.SessionTimeout = MIN_SESSION_TIMEOUT
	}
	if cfg.Cluster == "" {
		cfg.Cluster = DefaultClusterName
	}
	return nil
}

func (cfg *Proxy) ToBytes() ([]byte, error) {
	return yaml.Marshal(cfg)
}
