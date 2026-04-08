# Quick Start Guide

Get your Kostal Plenticore exporter up and running in minutes!

## 1. Prerequisites

- Your Kostal Plenticore inverter IP address
- Your inverter password
- Docker and Docker Compose installed (or Go 1.24+ for building from source)

## 2. Quick Setup with Docker

```bash
# 1. Clone or download this repository
git clone <your-repo-url>
cd kostal-plenticore-exporter

# 2. Create your configuration
cp .env.example .env

# 3. Edit .env with your inverter details
nano .env
# Update these values:
#   INVERTER_SERVER=192.168.1.100  (your inverter IP)
#   INVERTER_PASSWORD=your_password

# 4. Start the exporter
docker-compose up -d

# 5. Check if it's working
curl http://localhost:9876/health
curl http://localhost:9876/metrics
```

That's it! Your metrics are now available at `http://localhost:9876/metrics`

## 3. Complete Monitoring Stack (Optional)

Want a full monitoring solution with Prometheus and Grafana?

```bash
# Use the full stack compose file
docker-compose -f docker-compose.full-stack.yml up -d
```

This starts:
- Kostal Exporter (port 9876)
- Prometheus (port 9090)
- Grafana (port 3000)

Access Grafana at `http://localhost:3000` (default: admin/admin)

## 4. Verify Your Setup

### Check the exporter
```bash
# Health check
curl http://localhost:9876/health

# Sample metrics
curl http://localhost:9876/metrics | grep kostal_plenticore_devices_local

# Check Prometheus targets
open http://localhost:9090/targets
```

### View real-time metrics
```bash
# Current PV production
curl -s http://localhost:9876/metrics | grep 'kostal_plenticore_devices_local_pv.*Power'

# Battery state of charge
curl -s http://localhost:9876/metrics | grep 'Soc'

# Home consumption
curl -s http://localhost:9876/metrics | grep 'HomeBat_P'
```

## 5. Troubleshooting

### Exporter won't start

```bash
# Check logs
docker-compose logs -f

# Common issues:
# - Wrong IP address: ping your inverter first
# - Wrong password: test it in the inverter web UI
# - Network issues: ensure inverter is reachable
```

### No metrics showing

```bash
# Test the exporter directly
curl http://localhost:9876/metrics

# If empty, check authentication
docker-compose logs | grep -i "login\|auth\|error"
```

### Partial metrics

Different Kostal Plenticore models expose different metrics. To see what's available:

```bash
# List all available metrics
curl -s http://localhost:9876/metrics | grep kostal_plenticore | head -20
```

## 6. Next Steps

1. Configure Prometheus to scrape the exporter (see `prometheus.yml`)
2. Import the Grafana dashboard (see `grafana/` directory)
3. Set up alerts for your solar system
4. Customize metric collection intervals

See the [README.md](README.md) for detailed configuration options.
