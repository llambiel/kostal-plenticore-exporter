package main

import (
	"flag"
	"log"
	"net/http"
	"os"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listenAddress = flag.String("web.listen-address", getEnv("LISTEN_ADDRESS", ":9876"),
		"Address to listen on for web interface and telemetry")
	metricsPath = flag.String("web.telemetry-path", getEnv("METRICS_PATH", "/metrics"),
		"Path under which to expose metrics")
	inverterServer = flag.String("inverter.server", getEnv("INVERTER_SERVER", ""),
		"Hostname or IP address of Kostal Plenticore inverter")
	inverterPassword = flag.String("inverter.password", getEnv("INVERTER_PASSWORD", ""),
		"Password for Kostal Plenticore inverter")
	inverterScheme = flag.String("inverter.scheme", getEnv("INVERTER_SCHEME", "http"),
		"Scheme to use for connection (http or https)")
)

func getEnv(key, defaultValue string) string {
	if value := os.Getenv(key); value != "" {
		return value
	}
	return defaultValue
}

func main() {
	flag.Parse()

	if *inverterServer == "" {
		log.Fatal("Inverter server address is required. Set -inverter.server flag or INVERTER_SERVER environment variable")
	}

	if *inverterPassword == "" {
		log.Fatal("Inverter password is required. Set -inverter.password flag or INVERTER_PASSWORD environment variable")
	}

	log.Printf("Starting Kostal Plenticore Prometheus Exporter")
	log.Printf("Connecting to inverter at %s://%s", *inverterScheme, *inverterServer)
	log.Printf("Metrics available at %s%s", *listenAddress, *metricsPath)

	// Create collector
	collector := NewKostalCollector(*inverterServer, *inverterPassword, *inverterScheme)

	// Register collector
	prometheus.MustRegister(collector)

	// Setup HTTP handlers
	http.Handle(*metricsPath, promhttp.Handler())
	http.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Content-Type", "text/html")
		w.Write([]byte(`<html>
			<head><title>Kostal Plenticore Exporter</title></head>
			<body>
			<h1>Kostal Plenticore Prometheus Exporter</h1>
			<p><a href='` + *metricsPath + `'>Metrics</a></p>
			<h2>Build Info</h2>
			<ul>
				<li>Inverter: ` + *inverterServer + `</li>
				<li>Scheme: ` + *inverterScheme + `</li>
			</ul>
			</body>
			</html>`))
	})

	http.HandleFunc("/health", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		w.Write([]byte("OK"))
	})

	// Start server
	log.Fatal(http.ListenAndServe(*listenAddress, nil))
}
