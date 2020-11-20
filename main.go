package main

import (
	"fmt"
	"math"
	"math/rand"
	"net/http"
	"os"
	"strings"
	"sync"
	"time"

	"log"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	mtx          sync.Mutex
	testingID    string
	metricCount  = 1
	mc           = newMetricCollector()
	promRegistry = prometheus.NewRegistry() // local Registry so we don't get Go metrics, etc.
)

type metricBatch struct {
	counter   prometheus.Counter
	gauge     prometheus.Gauge
	histogram prometheus.Histogram
	summary   prometheus.Summary
}

func main() {
	testingID = os.Getenv("INSTANCE_ID")
	port := ":" + strings.Split(os.Getenv("LISTEN_ADDRESS"), ":")[1]
	rand.Seed(time.Now().Unix())

	registerMetrics(metricCount)
	go updateMetrics()

	http.HandleFunc("/", healthCheckHandler)
	http.Handle("/metrics", promhttp.HandlerFor(promRegistry, promhttp.HandlerOpts{}))
	http.HandleFunc("/expected_metrics", retrieveExpectedMetrics)

	log.Fatal(http.ListenAndServe(port, nil))
}

func updateMetrics() {
	for {
		time.Sleep(time.Second * 30)
		mc.timestamp = float64(time.Now().UnixNano()) / 1000000000
		for idx := 0; idx < mc.metricCount; idx++ {
			mtx.Lock()
			mc.counters[idx].Add(rand.Float64())
			mc.gauges[idx].Add(rand.Float64())
			lowerBound := math.Mod(rand.Float64(), 1)
			increment := math.Mod(rand.Float64(), 0.05)
			for i := lowerBound; i < 1; i += increment {
				mc.histograms[idx].Observe(i)
				mc.summarys[idx].Observe(i)
			}
			mtx.Unlock()
		}
	}
}

func retrieveExpectedMetrics(w http.ResponseWriter, r *http.Request) {
	mtx.Lock()
	defer mtx.Unlock()

	metricsResponse := mc.convertMetricsToExportedMetrics()
	retrieveExpectedMetricsHelper(w, r, metricsResponse)
}

func registerMetrics(metricCount int) {
	mc.metricCount = metricCount
	for idx := 0; idx < metricCount; idx++ {
		counter := prometheus.NewCounter(
			prometheus.CounterOpts{
				Namespace: testingID,
				Name:      fmt.Sprintf("test_counter%v", idx),
				Help:      "This is my counter",
				// labels can be added like this
				// ConstLabels: prometheus.Labels{
				// 	"label1": "val1",
				// },
			})
		gauge := prometheus.NewGauge(
			prometheus.GaugeOpts{
				Namespace: testingID,
				Name:      fmt.Sprintf("test_gauge%v", idx),
				Help:      "This is my gauge",
			})
		histogram := prometheus.NewHistogram(
			prometheus.HistogramOpts{
				Namespace: testingID,
				Name:      fmt.Sprintf("test_histogram%v", idx),
				Help:      "This is my histogram",
				Buckets:   []float64{0.005, 0.1, 1},
			})
		summary := prometheus.NewSummary(
			prometheus.SummaryOpts{
				Namespace: testingID,
				Name:      fmt.Sprintf("test_summary%v", idx),
				Help:      "This is my summary",
				Objectives: map[float64]float64{
					0.1:  0.5,
					0.5:  0.5,
					0.99: 0.5,
				},
			})

		promRegistry.MustRegister(counter)
		promRegistry.MustRegister(gauge)
		promRegistry.MustRegister(histogram)
		promRegistry.MustRegister(summary)

		mc.counters = append(mc.counters, counter)
		mc.gauges = append(mc.gauges, gauge)
		mc.histograms = append(mc.histograms, histogram)
		mc.summarys = append(mc.summarys, summary)
	}
}

func healthCheckHandler(w http.ResponseWriter, r *http.Request) {
	fmt.Fprintf(w, "healthy")
}
