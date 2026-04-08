# Kostal Plenticore Prometheus Exporter

A modern, comprehensive Prometheus exporter for Kostal Plenticore solar inverters written in Go. This exporter collects all available metrics from your inverter and exports them in Prometheus format for monitoring and visualization.

## Features

- **Comprehensive Metrics**: Exports all process data from all inverter modules
- **Easy Configuration**: Configure via environment variables or command-line flags
- **Docker Support**: Ready-to-use Docker image and docker-compose setup
- **Health Checks**: Built-in health endpoint for monitoring
- **Efficient**: Single binary with minimal resource usage
- **Modern**: Built with latest Go version and actively maintained libraries

## Metrics Collected

The exporter collects metrics from all Kostal Plenticore modules:

- **Energy Flow Statistics** (68 metrics): Daily, monthly, yearly, and lifetime energy statistics
  - Autarky rates
  - Self-consumption
  - CO2 savings
  - Battery charge/discharge
  - Grid import/export
  - PV generation
  - Home consumption

- **Real-time Device Data** (27 metrics): Instantaneous power and energy flows
- **AC Grid Monitoring** (16 metrics): Three-phase voltage, current, power, reactive power, frequency
- **Battery Status** (11 metrics): State of charge, cycles, capacity, temperature
- **Power Meter** (5 metrics): Grid consumption and feed-in
- **PV Panels** (6 metrics): Individual panel voltage, current, and power (PV1 & PV2)
- **System Events** (4 metrics): Error and warning counters

All metrics are exported with labels for module ID, metric ID, and unit for easy filtering and aggregation.

## Installation

### Using Docker (Recommended)

1. Clone this repository or create a new directory
2. Copy `.env.example` to `.env` and configure your inverter details:

```bash
cp .env.example .env
```

Edit `.env`:
```env
INVERTER_SERVER=192.168.1.100
INVERTER_PASSWORD=your_password
INVERTER_SCHEME=http
```

3. Start the exporter:

```bash
docker-compose up -d
```

The metrics will be available at `http://localhost:9876/metrics`

### Using Pre-built Binary

Download the latest release from the releases page and run:

```bash
./kostal-plenticore-exporter \
  -inverter.server=192.168.1.100 \
  -inverter.password=your_password \
  -inverter.scheme=http
```

### Building from Source

Requirements:
- Go 1.26 or later

```bash
# Clone the repository
git clone https://github.com/yourusername/kostal-plenticore-exporter
cd kostal-plenticore-exporter

# Download dependencies
go mod download

# Build
go build -o kostal-plenticore-exporter

# Run
./kostal-plenticore-exporter \
  -inverter.server=192.168.1.100 \
  -inverter.password=your_password
```

## Configuration

The exporter can be configured using environment variables or command-line flags.

### Environment Variables

| Variable | Description | Default |
|----------|-------------|---------|
| `INVERTER_SERVER` | Hostname or IP of your Kostal inverter | (required) |
| `INVERTER_PASSWORD` | Inverter password | (required) |
| `INVERTER_SCHEME` | Connection scheme (http or https) | `http` |
| `LISTEN_ADDRESS` | Address to listen on | `:9876` |
| `METRICS_PATH` | Path to expose metrics | `/metrics` |

### Command-Line Flags

| Flag | Description | Default |
|------|-------------|---------|
| `-inverter.server` | Hostname or IP of your Kostal inverter | (required) |
| `-inverter.password` | Inverter password | (required) |
| `-inverter.scheme` | Connection scheme (http or https) | `http` |
| `-web.listen-address` | Address to listen on | `:9876` |
| `-web.telemetry-path` | Path to expose metrics | `/metrics` |

## Prometheus Configuration

Add the exporter to your `prometheus.yml`:

```yaml
scrape_configs:
  - job_name: 'kostal-plenticore'
    static_configs:
      - targets: ['localhost:9876']
    scrape_interval: 30s
    scrape_timeout: 10s
```

### Example Prometheus Queries

```promql
# Current PV generation (Watts)
kostal_plenticore_devices_local_pv1_P{metric_id="Dc.1.Power"}

# Battery state of charge (%)
kostal_plenticore_devices_local_battery_SoC{metric_id="Soc"}

# Current home consumption (Watts)
kostal_plenticore_devices_local_HomeBat_P{metric_id="HomeBat_P"}

# Grid power (positive = import, negative = export)
kostal_plenticore_devices_local_ac_P{metric_id="P"}

# Total PV generation today (kWh)
kostal_plenticore_scb_statistic_energyflow_EnergyHomePv{metric_id="Statistic:Yield:Day"}
```

## Grafana Dashboard

A pre-configured Grafana dashboard is available in the `grafana/` directory. Import it to get started quickly with visualizations for:

- Real-time power flows
- Energy production and consumption
- Battery status
- Grid statistics
- Self-sufficiency metrics

## Endpoints

- `/metrics` - Prometheus metrics endpoint
- `/health` - Health check endpoint (returns 200 OK)
- `/` - Basic web interface with exporter info

## Troubleshooting

### Cannot connect to inverter

1. Verify the inverter IP address is correct
2. Check that you can ping the inverter from the host running the exporter
3. Ensure the password is correct (test it by logging into the inverter web interface)
4. If using HTTPS, make sure your inverter supports it

### No metrics showing up

1. Check the exporter logs: `docker-compose logs -f`
2. Verify the exporter is running: `curl http://localhost:9876/health`
3. Check if authentication is working by looking for login errors in logs
4. Ensure your Prometheus scrape configuration is correct

### Metrics seem incomplete

The exporter fetches all available metrics from your specific inverter model. Different Kostal Plenticore models may expose different sets of metrics. Check the `/metrics` endpoint directly to see what's available.

## Architecture

This exporter uses the [golrackpi](https://github.com/geschke/golrackpi) library to communicate with Kostal Plenticore inverters via their REST API. On each scrape:

1. The exporter logs into the inverter
2. Retrieves all available process data modules
3. Fetches metrics from each module
4. Converts values to Prometheus format
5. Logs out

The entire scrape process typically takes 1-3 seconds depending on your inverter and network.

## Credits

- Built using [golrackpi](https://github.com/geschke/golrackpi) by Ralf Geschke
- Inspired by the original [kostal-plenticore-exporter](https://codeberg.org/Marix/kostal-plenticore-exporter) Python implementation
- Vibe coded with [Claude Code](https://claude.com/claude-code)

## License

MIT License - see LICENSE file for details

## Contributing

Contributions are welcome! Please feel free to submit a Pull Request.

## Support

If you encounter any issues or have questions:

1. Check the [Troubleshooting](#troubleshooting) section
2. Review existing [GitHub Issues](https://github.com/yourusername/kostal-plenticore-exporter/issues)
3. Open a new issue with:
   - Your inverter model
   - Exporter version
   - Error logs
   - Configuration (without sensitive data)
