# IoTwin | Backend

[![GitHub go.mod Go version of a Go module](https://img.shields.io/github/go-mod/go-version/gomods/athens.svg)](https://github.com/muratalkan/iotwin-backend)
[![Go](https://github.com/muratalkan/iotwin-backend/actions/workflows/go.yml/badge.svg)](https://github.com/muratalkan/iotwin-backend/actions/workflows/go.yml)
[![CodeQL](https://github.com/muratalkan/iotwin-backend/actions/workflows/codeql.yml/badge.svg)](https://github.com/muratalkan/iotwin-backend/actions/workflows/codeql.yml)
[![MIT License](https://img.shields.io/badge/License-MIT-blue.svg)](LICENSE.md)


## Table of Contents
- [Features](#features)
- [Prerequisites](#prerequisites)
- [Setup](#setup)
- [Configuration](#configuration)
- [Roadmap](#roadmap)
- [License](#license)

-------

## Features
- Telemetry data collection (mqtt-subscriber)
- Telemetry data processing (kafka producer)
- Telemetry management (db crud)
- Alert management (db crud)

## Prerequisites
1. [Go](https://go.dev/dl/)
2. [Docker Engine](https://docs.docker.com/engine/install/)
2. [Docker Compose](https://docs.docker.com/compose/install/) (please see the [docker.compose](https://github.com/muratalkan/iotwin-backend/blob/master/docker/docker-compose.yml) file)
3. [Kafka](https://kafka.apache.org/downloads)
4. [MQTT Broker](https://mosquitto.org/download/)
5. DB (PostgreSQL/MySQL/CockroachDB)

## Setup

```bash
git clone https://github.com/muratalkan/iotwin-backend.git
go run ./iotwin-backend/main.go
```

## Configuration

#### default ".env" file

```ini
BACKEND_HOST=localhost 

DB_HOST=localhost
DB_PORT=
DB_NAME=
DB_USER=
DB_PASSWORD=
DB_CERT=
DB_KEY=

MQTT_HOST=localhost
MQTT_PORT=1883
MQTT_CLIENT=
MQTT_TOPIC=
MQTT_USER=
MQTT_PASSWORD=
MQTT_QOS=1
MQTT_CERT=
MQTT_KEY=

KAFKA_BROKERS=['locahost:9091', 'localhost:9092', 'localhost:9093']
```

## Roadmap
- Code Revision and Optimization
- Advanced Security
- Logging
- Dockerizing
- Testing
- ...

## License
[(Back to top)](#table-of-contents)

Licensed under the MIT License (MIT) 2022 - [Murat Alkan](https://github.com/muratalkan). Please have a look at the [LICENSE.md](LICENSE.md) for more details.
