package metrics

import (
	"fmt"
	"math"
	"math/rand"
	"time"

	"github.com/prometheus/client_golang/prometheus"
)

type metricCollector struct {
	counters   []prometheus.Counter
	gauges     []prometheus.Gauge
	histograms []prometheus.Histogram
	summarys   []prometheus.Summary
	interval   int
}

var (
	promRegistry = prometheus.NewRegistry() // local Registry so we don't get Go metrics, etc.
)

func newMetricCollector() metricCollector {
	return metricCollector{}
}

func (mc *metricCollector) updateCounter() {
	for {
		time.Sleep(time.Second * time.Duration(mc.interval))
		for idx := 0; idx < len(mc.counters); idx++ {
			mc.counters[idx].Add(rand.Float64())
		}
	}
}

func (mc *metricCollector) updateGauge() {
	for {
		time.Sleep(time.Second * time.Duration(mc.interval))
		for idx := 0; idx < len(mc.gauges); idx++ {
			mc.gauges[idx].Add(rand.Float64())
		}
	}
}

func (mc *metricCollector) updateHistogram() {
	for {
		time.Sleep(time.Second * time.Duration(mc.interval))

		for idx := 0; idx < len(mc.histograms); idx++ {
			lowerBound := math.Mod(rand.Float64(), 1)
			increment := math.Mod(rand.Float64(), 0.05)
			for i := lowerBound; i < 1; i += increment {
				mc.histograms[idx].Observe(i)

			}
		}
	}
}
func (mc *metricCollector) updateSummary() {
	for {
		time.Sleep(time.Second * time.Duration(mc.interval))

		for idx := 0; idx < len(mc.summarys); idx++ {
			lowerBound := math.Mod(rand.Float64(), 1)
			increment := math.Mod(rand.Float64(), 0.05)
			for i := lowerBound; i < 1; i += increment {
				mc.summarys[idx].Observe(i)

			}
		}
	}
}

func (mc *metricCollector) updateMetrics(count int) {
	for {
		time.Sleep(time.Second * time.Duration(mc.interval))
		for idx := 0; idx < count; idx++ {
			if mc.counters != nil {
				mc.counters[idx].Add(rand.Float64())
			}
			if mc.gauges != nil {
				mc.gauges[idx].Add(rand.Float64())
			}
			lowerBound := math.Mod(rand.Float64(), 1)
			increment := math.Mod(rand.Float64(), 0.05)
			for i := lowerBound; i < 1; i += increment {
				if mc.histograms != nil {
					mc.histograms[idx].Observe(i)
				}
				if mc.summarys != nil {
					mc.summarys[idx].Observe(i)
				}
			}
		}
	}
}

func (mc *metricCollector) registerCounter(count int) {
	for idx := 0; idx < count; idx++ {
		namespace := "test"
		counter := prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: namespace,
				Name:      fmt.Sprintf("counter%v", idx),
				Help:      "This is my counter",
			})

		promRegistry.MustRegister(counter)
		mc.counters = append(mc.counters, counter)
	}
}

func (mc *metricCollector) registerGauge(count int) {
	for idx := 0; idx < count; idx++ {
		namespace := "test"
		gauge := prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: namespace,
				Name:      fmt.Sprintf("gauge%v", idx),
				Help:      "This is my gauge",
			})
		promRegistry.MustRegister(gauge)
		mc.gauges = append(mc.gauges, gauge)
	}
}

func (mc *metricCollector) registerHistogram(count int) {
	for idx := 0; idx < count; idx++ {
		namespace := "test"
		histogram := prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: namespace,
				Name:      fmt.Sprintf("histogram%v", idx),
				Help:      "This is my histogram",
				Buckets:   []float64{0.1, 0.5, 1},
			})
		promRegistry.MustRegister(histogram)
		mc.histograms = append(mc.histograms, histogram)
	}
}

func (mc *metricCollector) registerSummary(count int) {
	for idx := 0; idx < count; idx++ {
		namespace := "test"
		summary := prometheus.NewSummary(
			prometheus.SummaryOpts{
				Namespace: namespace,
				Name:      fmt.Sprintf("summary%v", idx),
				Help:      "This is my summary",
				Objectives: map[float64]float64{
					0.1:  0.5,
					0.5:  0.5,
					0.99: 0.5,
				},
			})
		promRegistry.MustRegister(summary)
		mc.summarys = append(mc.summarys, summary)
	}
}
