package metrics

import (
	"os"

	"github.com/prometheus/client_golang/prometheus"
)

const (
	MetricsNamespace           = "focalboard"
	MetricsSubsystemBlocks     = "blocks"
	MetricsSubsystemWorkspaces = "workspaces"
	MetricsSubsystemSystem     = "system"

	MetricsCloudInstallationLabel = "installationId"
)

type InstanceInfo struct {
	Version        string
	BuildNum       string
	Edition        string
	InstallationID string
}

// Metrics used to instrumentate metrics in prometheus
type Metrics struct {
	registry *prometheus.Registry

	instance  *prometheus.GaugeVec
	startTime prometheus.Gauge

	loginCount     prometheus.Counter
	loginFailCount prometheus.Counter

	blocksInsertedCount prometheus.Counter
	blocksDeletedCount  prometheus.Counter

	blockCount     *prometheus.GaugeVec
	workspaceCount prometheus.Gauge

	blockLastActivity prometheus.Gauge
}

// NewMetrics Factory method to create a new metrics collector
func NewMetrics(info InstanceInfo) *Metrics {
	m := &Metrics{}

	m.registry = prometheus.NewRegistry()
	options := prometheus.ProcessCollectorOpts{
		Namespace: MetricsNamespace,
	}
	m.registry.MustRegister(prometheus.NewProcessCollector(options))
	m.registry.MustRegister(prometheus.NewGoCollector())

	additionalLabels := map[string]string{}
	if info.InstallationID != "" {
		additionalLabels[MetricsCloudInstallationLabel] = os.Getenv("MM_CLOUD_INSTALLATION_ID")
	}

	m.loginCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   MetricsNamespace,
		Subsystem:   MetricsSubsystemSystem,
		Name:        "login_total",
		Help:        "Total number of logins.",
		ConstLabels: additionalLabels,
	})
	m.registry.MustRegister(m.loginCount)

	m.loginFailCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   MetricsNamespace,
		Subsystem:   MetricsSubsystemSystem,
		Name:        "login_fail_total",
		Help:        "Total number of failed logins.",
		ConstLabels: additionalLabels,
	})
	m.registry.MustRegister(m.loginFailCount)

	m.instance = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   MetricsNamespace,
		Subsystem:   MetricsSubsystemSystem,
		Name:        "focalboard_instance_info",
		Help:        "Instance information for Focalboard.",
		ConstLabels: additionalLabels,
	}, []string{"Version", "BuildNum", "Edition"})
	m.registry.MustRegister(m.instance)
	m.instance.WithLabelValues(info.Version, info.BuildNum, info.Edition).Set(1)

	m.startTime = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   MetricsNamespace,
		Subsystem:   MetricsSubsystemSystem,
		Name:        "server_start_time",
		Help:        "The time the server started.",
		ConstLabels: additionalLabels,
	})
	m.startTime.SetToCurrentTime()
	m.registry.MustRegister(m.startTime)

	m.blocksInsertedCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   MetricsNamespace,
		Subsystem:   MetricsSubsystemBlocks,
		Name:        "blocks_inserted_total",
		Help:        "Total number of blocks inserted.",
		ConstLabels: additionalLabels,
	})
	m.registry.MustRegister(m.blocksInsertedCount)

	m.blocksDeletedCount = prometheus.NewCounter(prometheus.CounterOpts{
		Namespace:   MetricsNamespace,
		Subsystem:   MetricsSubsystemBlocks,
		Name:        "blocks_deleted_total",
		Help:        "Total number of blocks deleted.",
		ConstLabels: additionalLabels,
	})
	m.registry.MustRegister(m.blocksDeletedCount)

	m.blockCount = prometheus.NewGaugeVec(prometheus.GaugeOpts{
		Namespace:   MetricsNamespace,
		Subsystem:   MetricsSubsystemBlocks,
		Name:        "blocks_total",
		Help:        "Total number of blocks.",
		ConstLabels: additionalLabels,
	}, []string{"BlockType"})
	m.registry.MustRegister(m.blockCount)

	m.workspaceCount = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   MetricsNamespace,
		Subsystem:   MetricsSubsystemWorkspaces,
		Name:        "workspaces_total",
		Help:        "Total number of workspaces.",
		ConstLabels: additionalLabels,
	})
	m.registry.MustRegister(m.workspaceCount)

	m.blockLastActivity = prometheus.NewGauge(prometheus.GaugeOpts{
		Namespace:   MetricsNamespace,
		Subsystem:   MetricsSubsystemBlocks,
		Name:        "blocks_last_activity",
		Help:        "Time of last block insert, update, delete.",
		ConstLabels: additionalLabels,
	})
	m.registry.MustRegister(m.blockLastActivity)

	return m
}

func (m *Metrics) IncrementLoginCount(num int) {
	if m != nil {
		m.loginCount.Add(float64(num))
	}
}

func (m *Metrics) IncrementLoginFailCount(num int) {
	if m != nil {
		m.loginFailCount.Add(float64(num))
	}
}

func (m *Metrics) IncrementBlocksInserted(num int) {
	if m != nil {
		m.blocksInsertedCount.Add(float64(num))
		m.blockLastActivity.SetToCurrentTime()
	}
}

func (m *Metrics) IncrementBlocksDeleted(num int) {
	if m != nil {
		m.blocksDeletedCount.Add(float64(num))
		m.blockLastActivity.SetToCurrentTime()
	}
}

func (m *Metrics) ObserveBlockCount(blockType string, count int64) {
	if m != nil {
		m.blockCount.WithLabelValues(blockType).Set(float64(count))
	}
}

func (m *Metrics) ObserveWorkspaceCount(count int64) {
	if m != nil {
		m.workspaceCount.Set(float64(count))
	}
}
