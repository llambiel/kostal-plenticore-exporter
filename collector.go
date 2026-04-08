package main

import (
	"fmt"
	"log"
	"strconv"
	"strings"
	"sync"

	"github.com/geschke/golrackpi"
	"github.com/prometheus/client_golang/prometheus"
)

// KostalCollector implements prometheus.Collector
type KostalCollector struct {
	client      *golrackpi.AuthClient
	mu          sync.Mutex
	up          *prometheus.Desc
	scrapeDuration *prometheus.Desc
	totalScrapes   prometheus.Counter
}

// NewKostalCollector creates a new collector
func NewKostalCollector(server, password, scheme string) *KostalCollector {
	client := golrackpi.NewWithParameter(golrackpi.AuthClient{
		Server:   server,
		Password: password,
		Scheme:   scheme,
	})

	return &KostalCollector{
		client: client,
		up: prometheus.NewDesc(
			"kostal_plenticore_up",
			"Was the last scrape of the Kostal Plenticore successful",
			nil, nil,
		),
		scrapeDuration: prometheus.NewDesc(
			"kostal_plenticore_scrape_duration_seconds",
			"Duration of the last scrape in seconds",
			nil, nil,
		),
		totalScrapes: prometheus.NewCounter(prometheus.CounterOpts{
			Name: "kostal_plenticore_scrapes_total",
			Help: "Total number of times Kostal Plenticore was scraped for metrics",
		}),
	}
}

// Describe sends the descriptors to the provided channel
func (c *KostalCollector) Describe(ch chan<- *prometheus.Desc) {
	ch <- c.up
	ch <- c.scrapeDuration
	c.totalScrapes.Describe(ch)
}

// Collect fetches metrics from the Kostal inverter
func (c *KostalCollector) Collect(ch chan<- prometheus.Metric) {
	c.mu.Lock()
	defer c.mu.Unlock()

	c.totalScrapes.Inc()
	c.totalScrapes.Collect(ch)

	start := prometheus.NewTimer(prometheus.ObserverFunc(func(v float64) {
		ch <- prometheus.MustNewConstMetric(c.scrapeDuration, prometheus.GaugeValue, v)
	}))
	defer start.ObserveDuration()

	// Login to the inverter
	_, err := c.client.Login()
	if err != nil {
		log.Printf("Failed to login: %v", err)
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		return
	}
	defer c.client.Logout()

	// Get all process data
	processData, err := c.client.ProcessData()
	if err != nil {
		log.Printf("Failed to get process data: %v", err)
		ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 0)
		return
	}

	// Collect metrics from all modules
	for _, module := range processData {
		c.collectModuleMetrics(ch, module.ModuleId)
	}

	ch <- prometheus.MustNewConstMetric(c.up, prometheus.GaugeValue, 1)
}

// collectModuleMetrics fetches and exports metrics for a specific module
func (c *KostalCollector) collectModuleMetrics(ch chan<- prometheus.Metric, moduleId string) {
	values, err := c.client.ProcessDataModule(moduleId)
	if err != nil {
		log.Printf("Failed to get process data for module %s: %v", moduleId, err)
		return
	}

	for _, moduleData := range values {
		for _, processData := range moduleData.ProcessData {
			// Convert the value to float64
			floatValue, err := convertToFloat64(processData.Value)
			if err != nil {
				log.Printf("Failed to convert value for %s/%s: %v", moduleId, processData.Id, err)
				continue
			}

			// Create metric name and labels
			metricName := sanitizeMetricName(fmt.Sprintf("kostal_plenticore_%s_%s", moduleId, processData.Id))

			labels := []string{"module_id", "metric_id", "unit"}
			labelValues := []string{moduleId, processData.Id, processData.Unit}

			desc := prometheus.NewDesc(
				metricName,
				fmt.Sprintf("Kostal Plenticore metric %s from module %s", processData.Id, moduleId),
				labels,
				nil,
			)

			ch <- prometheus.MustNewConstMetric(desc, prometheus.GaugeValue, floatValue, labelValues...)
		}
	}
}

// sanitizeMetricName converts a string into a valid Prometheus metric name
func sanitizeMetricName(name string) string {
	// Replace invalid characters with underscores
	name = strings.ReplaceAll(name, ":", "_")
	name = strings.ReplaceAll(name, "-", "_")
	name = strings.ReplaceAll(name, ".", "_")
	name = strings.ReplaceAll(name, "/", "_")
	name = strings.ToLower(name)
	return name
}

// convertToFloat64 attempts to convert various types to float64
func convertToFloat64(value interface{}) (float64, error) {
	switch v := value.(type) {
	case float64:
		return v, nil
	case float32:
		return float64(v), nil
	case int:
		return float64(v), nil
	case int32:
		return float64(v), nil
	case int64:
		return float64(v), nil
	case uint:
		return float64(v), nil
	case uint32:
		return float64(v), nil
	case uint64:
		return float64(v), nil
	case string:
		// Try to parse string as float
		f, err := strconv.ParseFloat(v, 64)
		if err != nil {
			// Try to parse as boolean
			if v == "true" || v == "on" {
				return 1, nil
			} else if v == "false" || v == "off" {
				return 0, nil
			}
			return 0, fmt.Errorf("cannot convert string to float64: %s", v)
		}
		return f, nil
	case bool:
		if v {
			return 1, nil
		}
		return 0, nil
	default:
		return 0, fmt.Errorf("unsupported type: %T", v)
	}
}
