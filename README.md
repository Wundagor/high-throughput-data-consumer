# High Throughput Data Consumer

A simple and scalable Go application that consumes messages from a RabbitMQ queue and stores them in a PostgreSQL database. This application is designed for high throughput and can handle a large volume of messages efficiently.

## Features

- Connects to RabbitMQ and consumes messages from a specified queue.
- Processes messages and saves them into a MySQL database.
- Implements a worker pool for concurrent message processing.
- Supports graceful shutdown and error handling.
- Easy to extend and modify for future requirements.

## Table of Contents

- [Requirements](#requirements)
- [Installation](#installation)

## Requirements

- [high-throughput-data-transfer](https://github.com/Wundagor/high-throughput-data-transfer)

## Installation

1. Clone the repository:
   ```bash
   git clone https://github.com/yourusername/high-throughput-data-consumer.git
   cd high-throughput-data-consumer

2. Run the docker:
   ```bash
   docker-compose up -d
