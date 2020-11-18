package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"sync"
	"time"

	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mtx             sync.Mutex
	testingID       string
	metricCount     = 1
	metricCollector MetricCollector
	promRegistry    = prometheus.NewRegistry() // local Registry so we don't get Go metrics, etc.
)

type metricBatch struct {
	counter   prometheus.Counter
	gauge     prometheus.Gauge
	histogram prometheus.Histogram
	summary   prometheus.Summary
}

func main() {
	rand.Seed(time.Now().Unix())

	metrics := registerMetrics(metricCount)
	go updateMetrics(metrics)

	http.HandleFunc("/", healthCheckHandler)
	http.Handle("/metrics", promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}))
	http.HandleFunc("/expected_metrics", retrieveExpectedMetrics)

	log.Fatal(http.ListenAndServe(":8080", nil))
}

func updateMetrics(metrics []*metricBatch) {
	for {
		for idx := range metrics {
			metrics[idx].counter.Add(rand.Float64())
			metrics[idx].gauge.Add(rand.Float64())
			lowerBound := math.Mod(rand.Float64(), 1)
			increment := math.Mod(rand.Float64(), 0.05)
			for i := lowerBound; i < 1; i += increment {
				metrics[idx].histogram.Observe(i)
				metrics[idx].summary.Observe(i)
			}
		}

		time.Sleep(time.Second)
	}
}

func retrieveExpectedMetrics(w http.ResponseWriter, r *http.Request) {
	mtx.Lock()
	defer mtx.Unlock()

	metricsResponse := metricCollector.convertMetricsToExportedMetrics()
	retrieveExpectedMetricsHelper(w, r, metricsResponse)
}

func registerMetrics(metricCount int) []*metricBatch {
	metrics := make([]*metricBatch, metricCount)
	for idx := 0; idx < metricCount; idx++ {
		counter := prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: testingID,
				Name:      fmt.Sprintf("test_counter_%v", idx),
				Help:      "This is my counter",
			})
		gauge := prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: testingID,
				Name:      fmt.Sprintf("test_gauge_%v", idx),
				Help:      "This is my gauge",
			})
		histogram := prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: testingID,
				Name:      fmt.Sprintf("test_histogram_%v", idx),
				Help:      "This is my histogram",
				Buckets:   []float64{0.005, 0.1, 1},
			})
		summary := prometheus.NewSummary(
			prometheus.SummaryOpts{
				Namespace: testingID,
				Name:      fmt.Sprintf("test_summary_%v", idx),
				Help:      "This is my summary",
				Objectives: map[float64]float64{
					0.1:  0.5,
					0.5:  0.5,
					0.99: 0.5,
				},
			})

		// Set arbitrary values
		counter.Add(5)
		gauge.Add(15)
		for i := 0.005; i < 1; i += 0.005 {
			histogram.Observe(i)
			summary.Observe(i)
		}

		promRegistry.MustRegister(counter)
		promRegistry.MustRegister(gauge)
		promRegistry.MustRegister(histogram)
		promRegistry.MustRegister(summary)

		newMetrics := metricBatch{counter: counter, gauge: gauge, histogram: histogram, summary: summary}
		metrics[idx] = &newMetrics
	}
	return metrics
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "healthy")
}
