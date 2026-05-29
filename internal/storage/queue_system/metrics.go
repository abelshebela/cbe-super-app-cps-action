package queue

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// QueueMetrics holds Prometheus metrics for queue operations.
type QueueMetrics struct {
	enqueued     *prometheus.CounterVec
	processed    *prometheus.CounterVec
	failed       *prometheus.CounterVec
	retried      *prometheus.CounterVec
	deadLettered *prometheus.CounterVec
	deduplicated *prometheus.CounterVec
	depth        *prometheus.GaugeVec
}

// NewQueueMetrics registers and returns a new set of queue metrics.
// The namespace parameter (e.g. "cps_action") prefixes all metric names.
func NewQueueMetrics(namespace string) *QueueMetrics {
	return &QueueMetrics{
		enqueued: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "queue",
			Name:      "jobs_enqueued_total",
			Help:      "Total number of jobs enqueued.",
		}, []string{"backend", "job_type"}),

		processed: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "queue",
			Name:      "jobs_processed_total",
			Help:      "Total number of jobs successfully processed.",
		}, []string{"backend", "job_type"}),

		failed: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "queue",
			Name:      "jobs_failed_total",
			Help:      "Total number of jobs that failed processing.",
		}, []string{"backend", "job_type"}),

		retried: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "queue",
			Name:      "jobs_retried_total",
			Help:      "Total number of jobs scheduled for retry.",
		}, []string{"backend", "job_type"}),

		deadLettered: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "queue",
			Name:      "jobs_dead_lettered_total",
			Help:      "Total number of jobs moved to the dead letter queue.",
		}, []string{"backend", "job_type"}),

		deduplicated: promauto.NewCounterVec(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: "queue",
			Name:      "jobs_deduplicated_total",
			Help:      "Total number of jobs skipped due to deduplication.",
		}, []string{"job_type"}),

		depth: promauto.NewGaugeVec(prometheus.GaugeOpts{
			Namespace: namespace,
			Subsystem: "queue",
			Name:      "depth",
			Help:      "Current number of jobs waiting in the queue.",
		}, []string{"backend"}),
	}
}

// --- Counter helpers ---

func (m *QueueMetrics) IncEnqueued(backend, jobType string) {
	if m != nil {
		m.enqueued.WithLabelValues(backend, jobType).Inc()
	}
}

func (m *QueueMetrics) IncProcessed(backend, jobType string) {
	if m != nil {
		m.processed.WithLabelValues(backend, jobType).Inc()
	}
}

func (m *QueueMetrics) IncFailed(backend, jobType string) {
	if m != nil {
		m.failed.WithLabelValues(backend, jobType).Inc()
	}
}

func (m *QueueMetrics) IncRetried(backend, jobType string) {
	if m != nil {
		m.retried.WithLabelValues(backend, jobType).Inc()
	}
}

func (m *QueueMetrics) IncDeadLettered(backend, jobType string) {
	if m != nil {
		m.deadLettered.WithLabelValues(backend, jobType).Inc()
	}
}

func (m *QueueMetrics) IncDeduplicated(jobType string) {
	if m != nil {
		m.deduplicated.WithLabelValues(jobType).Inc()
	}
}

// --- Gauge helpers ---

func (m *QueueMetrics) SetDepth(backend string, value float64) {
	if m != nil {
		m.depth.WithLabelValues(backend).Set(value)
	}
}
