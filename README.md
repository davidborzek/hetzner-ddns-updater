[![Tests](https://github.com/davidborzek/hetzner-ddns-updater/actions/workflows/tests.yml/badge.svg)](https://github.com/davidborzek/hetzner-ddns-updater/actions/workflows/tests.yml)

# hetzner-ddns-updater

hetzner-ddns-updater is a lightweight service designed to periodically check for changes in your external IP address and update a DNS record at Hetzner when necessary. This project simplifies the process of keeping your DNS records up-to-date with your dynamic IP address.

## Table of Contents

- [Configuration](#configuration)
- [Metrics](#metrics)
- [Running with Docker](#running-with-docker)

## Configuration

To configure hetzner-ddns-updater, you can use environment variables. Here are the available configuration parameters:

| Environment Variable  | Default Value | Description                                                      |
| --------------------- | ------------- | ---------------------------------------------------------------- |
| `HDU_ADDRESS`         | `:8080`       | The address and port on which the service will listen.           |
| `HDU_API_TOKEN`       | Required      | Your Hetzner API token for authentication.                       |
| `HDU_RECORD_ID`       | Required      | The ID of the DNS record to be updated.                          |
| `HDU_ZONE_ID`         | Required      | The ID of the DNS zone where the record resides.                 |
| `HDU_RECORD_NAME`     | `@`           | The DNS record name to be updated (e.g., subdomain).             |
| `HDU_RECORD_TTL`      | `60`          | Time to live (TTL) for the DNS record in seconds.                |
| `HDU_INTERVAL`        | `5m`          | The interval at which the service checks for IP address changes. |
| `HDU_METRICS_ENABLED` | `false`       | Enable or disable Prometheus metrics.                            |
| `HDU_METRICS_TOKEN`   |               | Token to secure access to Prometheus metrics when enabled.       |

## Metrics

When Prometheus metrics are enabled (`HDU_METRICS_ENABLED=true`), the service exposes the following metrics in Prometheus format at the `/metrics` route:

| Name               | Type    | Description                   |
| ------------------ | ------- | ----------------------------- |
| hdu_total_updates  | Counter | Number of ddns updates        |
| hdu_failed_updates | Counter | Number of failed ddns updates |

These metrics provide insights into the update activity of the service, helping you monitor its performance and reliability.

## Running with Docker

You can easily run hetzner-ddns-updater as a Docker container. Here's an example command to start the service:

```shell
docker run -d \
     -e "HDU_API_TOKEN=<your_hetzner_api_token>" \
     -e "HDU_RECORD_ID=<your_record_id>" \
     -e "HDU_ZONE_ID=<your_zone_id>" \
     -p 8080:8080 \
     ghcr.io/davidborzek/hetzner-ddns-updater:latest
```
