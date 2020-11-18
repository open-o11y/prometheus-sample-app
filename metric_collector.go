package main

import (
	"fmt"

	"github.com/prometheus/client_golang/prometheus"
	dto "github.com/prometheus/client_model/go"
)

type MetricCollector struct {
	counters   []prometheus.Counter
	gauges     []prometheus.Gauge
	histograms []prometheus.Histogram
	summarys   []prometheus.Summary

	metricCount int
	timestamp   float64
}

func (mc MetricCollector) convertMetricsToExportedMetrics() []MetricResponse {
	metricsResponse := make([]MetricResponse, metricCollector.metricCount)

	mc.handleCounters(metricsResponse)
	mc.handleGauges(metricsResponse)
	mc.handleHistograms(metricsResponse)
	mc.handleSummarys(metricsResponse)

	return metricsResponse
}

func (mc MetricCollector) handleCounters(metricsResponse []MetricResponse) {
	for _, counter := range metricCollector.counters {
		var metric *dto.Metric
		counter.Write(metric)

		labels := convertLabelPairsToLabels(metric.GetLabel())
		values := convertMetricValues(metricCollector.timestamp, metric.GetCounter().GetValue())

		metricsResponse = append(metricsResponse, MetricResponse{
			Labels: labels,
			Value:  values,
		})
	}
}

func (mc MetricCollector) handleGauges(metricsResponse []MetricResponse) {
	for _, gauge := range metricCollector.gauges {
		var metric *dto.Metric
		gauge.Write(metric)

		labels := convertLabelPairsToLabels(metric.GetLabel())
		values := convertMetricValues(metricCollector.timestamp, metric.GetGauge().GetValue())

		metricsResponse = append(metricsResponse, MetricResponse{
			Labels: labels,
			Value:  values,
		})
	}
}

func (mc MetricCollector) handleHistograms(metricsResponse []MetricResponse) {
	for _, histogram := range metricCollector.histograms {
		var metric *dto.Metric
		histogram.Write(metric)

		// handle count
		countLabels := convertLabelPairsToLabels(metric.GetLabel())
		countLabels["__name__"] += "_count"
		countValues := convertMetricValues(metricCollector.timestamp, float64(metric.GetHistogram().GetSampleCount()))

		metricsResponse = append(metricsResponse, MetricResponse{
			Labels: countLabels,
			Value:  countValues,
		})

		// handle sum
		sumLabels := convertLabelPairsToLabels(metric.GetLabel())
		sumLabels["__name__"] += "_sum"
		sumValues := convertMetricValues(metricCollector.timestamp, metric.GetHistogram().GetSampleSum())

		metricsResponse = append(metricsResponse, MetricResponse{
			Labels: sumLabels,
			Value:  sumValues,
		})

		// handle buckets
		for _, bucket := range metric.GetHistogram().GetBucket() {
			labels := convertLabelPairsToLabels(metric.GetLabel())
			labels["__name__"] += "_bucket"
			labels["le"] = fmt.Sprintf("%f", bucket.GetUpperBound())
			values := convertMetricValues(metricCollector.timestamp, float64(bucket.GetCumulativeCount()))

			metricsResponse = append(metricsResponse, MetricResponse{
				Labels: labels,
				Value:  values,
			})
		}
	}
}

func (mc MetricCollector) handleSummarys(metricsResponse []MetricResponse) {
	for _, summary := range metricCollector.summarys {
		var metric *dto.Metric
		summary.Write(metric)

		// handle count
		countLabels := convertLabelPairsToLabels(metric.GetLabel())
		countLabels["__name__"] += "_count"
		countValues := convertMetricValues(metricCollector.timestamp, float64(metric.GetSummary().GetSampleCount()))

		metricsResponse = append(metricsResponse, MetricResponse{
			Labels: countLabels,
			Value:  countValues,
		})

		// handle sum
		sumLabels := convertLabelPairsToLabels(metric.GetLabel())
		sumLabels["__name__"] += "_sum"
		sumValues := convertMetricValues(metricCollector.timestamp, metric.GetSummary().GetSampleSum())

		metricsResponse = append(metricsResponse, MetricResponse{
			Labels: sumLabels,
			Value:  sumValues,
		})

		// handle quantiles
		for _, quantile := range metric.GetSummary().GetQuantile() {
			labels := convertLabelPairsToLabels(metric.GetLabel())
			labels["quantile"] = fmt.Sprintf("%f", quantile.GetQuantile())
			values := convertMetricValues(metricCollector.timestamp, float64(quantile.GetValue()))

			metricsResponse = append(metricsResponse, MetricResponse{
				Labels: labels,
				Value:  values,
			})
		}
	}
}

func convertLabelPairsToLabels(labelPairs []*dto.LabelPair) map[string]string {
	var labels map[string]string
	for _, labelPair := range labelPairs {
		labels[labelPair.GetName()] = labelPair.GetValue()
	}
	return labels
}

func convertMetricValues(timestamp float64, value float64) []string {
	values := make([]string, 2)
	values[0] = fmt.Sprintf("%f", timestamp)
	values[1] = fmt.Sprintf("%f", value)
	return values
}
