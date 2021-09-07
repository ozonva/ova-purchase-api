package metrics

import (
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promauto"
)

// Metrics - interface for Create/Update/Remove operations on Purchase
type Metrics interface {
	CreatePurchaseCounterInc()
	MultiCreatePurchaseCounterInc()
	UpdatePurchaseCounterInc()
	RemovePurchaseCounterInc()
}

type metrics struct {
	createPurchaseSuccessCounter      prometheus.Counter
	multiCreatePurchaseSuccessCounter prometheus.Counter
	updatePurchaseSuccessCounter      prometheus.Counter
	removePurchaseSuccessCounter      prometheus.Counter
}

// NewMetrics - creates new Metrics object
func NewMetrics(namespace, subsystem string) Metrics {
	return &metrics{
		createPurchaseSuccessCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "create_count_total",
			Help:      "Total count of successful requests to create Purchase",
		}),
		multiCreatePurchaseSuccessCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "multi_create_count_total",
			Help:      "Total count of successful requests to chunked create Purchases",
		}),
		updatePurchaseSuccessCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "update_count_total",
			Help:      "Total count of successful requests to update Purchase",
		}),
		removePurchaseSuccessCounter: promauto.NewCounter(prometheus.CounterOpts{
			Namespace: namespace,
			Subsystem: subsystem,
			Name:      "remove_count_total",
			Help:      "Total count of successful requests to remove Purchase",
		}),
	}
}

func (m *metrics) CreatePurchaseCounterInc() {
	m.createPurchaseSuccessCounter.Inc()
}

func (m *metrics) MultiCreatePurchaseCounterInc() {
	m.multiCreatePurchaseSuccessCounter.Inc()
}

func (m *metrics) UpdatePurchaseCounterInc() {
	m.updatePurchaseSuccessCounter.Inc()
}

func (m *metrics) RemovePurchaseCounterInc() {
	m.removePurchaseSuccessCounter.Inc()
}
