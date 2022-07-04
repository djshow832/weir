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

package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
)

type AstStmtType int

const (
	StmtTypeUnknown AstStmtType = iota
	StmtTypeSelect
	StmtTypeInsert
	StmtTypeUpdate
	StmtTypeDelete
	StmtTypeDDL
	StmtTypeBegin
	StmtTypeCommit
	StmtTypeRollback
	StmtTypeSet
	StmtTypeShow
	StmtTypeUse
	StmtTypeComment
)

const (
	StmtNameUnknown  = "unknown"
	StmtNameSelect   = "select"
	StmtNameInsert   = "insert"
	StmtNameUpdate   = "update"
	StmtNameDelete   = "delete"
	StmtNameDDL      = "ddl"
	StmtNameBegin    = "begin"
	StmtNameCommit   = "commit"
	StmtNameRollback = "rollback"
	StmtNameSet      = "set"
	StmtNameShow     = "show"
	StmtNameUse      = "use"
	StmtNameComment  = "comment"
)

var (
	QueryCtxQueryCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ModuleWeirProxy,
			Subsystem: LabelQueryCtx,
			Name:      "query_total",
			Help:      "Counter of queries.",
		}, []string{LblCluster, LblNamespace, LblDb, LblTable, LblSQLType, LblResult})

	QueryCtxQueryDeniedCounter = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Namespace: ModuleWeirProxy,
			Subsystem: LabelQueryCtx,
			Name:      "query_denied",
			Help:      "Counter of denied queries.",
		}, []string{LblCluster, LblNamespace, LblDb, LblTable, LblSQLType})

	QueryCtxQueryDurationHistogram = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: ModuleWeirProxy,
			Subsystem: LabelQueryCtx,
			Name:      "handle_query_duration_seconds",
			Help:      "Bucketed histogram of processing time (s) of handled queries.",
			Buckets:   prometheus.ExponentialBuckets(0.0005, 2, 29), // 0.5ms ~ 1.5days
		}, []string{LblCluster, LblNamespace, LblDb, LblTable, LblSQLType})

	QueryCtxGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: ModuleWeirProxy,
			Subsystem: LabelQueryCtx,
			Name:      "queryctx",
			Help:      "Number of queryctx (equals to client connection).",
		}, []string{LblCluster, LblNamespace})

	QueryCtxAttachedConnGauge = prometheus.NewGaugeVec(
		prometheus.GaugeOpts{
			Namespace: ModuleWeirProxy,
			Subsystem: LabelQueryCtx,
			Name:      "attached_connections",
			Help:      "Number of attached backend connections.",
		}, []string{LblCluster, LblNamespace})

	QueryCtxTransactionDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Namespace: "tidb",
			Subsystem: "session",
			Name:      "transaction_duration_seconds",
			Help:      "Bucketed histogram of a transaction execution duration, including retry.",
			Buckets:   prometheus.ExponentialBuckets(0.001, 2, 28), // 1ms ~ 1.5days
		}, []string{LblCluster, LblNamespace, LblDb, LblSQLType})
)
