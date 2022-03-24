package config

type Namespace struct {
	Version     string            `yaml:"version"`
	Namespace   string            `yaml:"namespace"`
	Frontend    FrontendNamespace `yaml:"frontend"`
	Backend     BackendNamespace  `yaml:"backend"`
	Breaker     BreakerInfo       `yaml:"breaker"`
	RateLimiter RateLimiterInfo   `yaml:"rate_limiter"`
}

type FrontendNamespace struct {
	AllowedDBs   []string  `yaml:"allowed_dbs"`
	SlowSQLTime  int       `yaml:"slow_sql_time"`
	DeniedIPs    []string  `yaml:"denied_ips"`
	Usernames    []string  `yaml:"usernames"`
	SQLBlackList []SQLInfo `yaml:"sql_blacklist"`
	SQLWhiteList []SQLInfo `yaml:"sql_whitelist"`
}

type SQLInfo struct {
	SQL string `yaml:"sql"`
}

type RateLimiterInfo struct {
	Scope string `yaml:"scope"`
	QPS   int    `yaml:"qps"`
}

type BackendNamespace struct {
	Instances    []string `yaml:"instances"`
	SelectorType string   `yaml:"selector_type"`
}

type StrategyInfo struct {
	MinQps               int64 `yaml:"min_qps"`
	SqlTimeoutMs         int64 `yaml:"sql_timeout_ms"`
	FailureRatethreshold int64 `yaml:"failure_rate_threshold"`
	FailureNum           int64 `yaml:"failure_num"`
	OpenStatusDurationMs int64 `yaml:"open_status_duration_ms"`
	Size                 int64 `yaml:"size"`
	CellIntervalMs       int64 `yaml:"cell_interval_ms"`
}

type BreakerInfo struct {
	Scope      string         `yaml:"scope"`
	Strategies []StrategyInfo `yaml:"strategies"`
}
