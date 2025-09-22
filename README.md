# R1 Cloud

---

## SMS Gateway Code Challenge

✅ Executed based on assessment guidelines:

- Used GO version: 1.24
- Used PostgreSQL Database
- Designed Hexagonal Architecture
- SMS Provider Mocked and Integrated
- Dockerized via Dockerfile and Docker Compose
- Initialized by Makefile and make command
- Send SMS to optional Mobile number by API
- List of sent messages(sent messages report)
- Increase credit to send message per tenant
- Shared received `Send Message` requests between `Kafka` topics
- `Normal` and `Express` SMS segregated
- Checks messages length count and price amount to avoid send with no credit
- The credit transactions list (increase credit balance/decrease credit per sent SMS)

---

## Overview

### Service:

1. First, create a tenant that has `create`, `detail`, and `list` APIs
2. The `Tenant` balance shown in `Detail` 
3. Use the Tenant `UUID` in all other requests' headers
4. Increase the `Credit` to send SMS(Use the `list` API to trace transactions)
5. Send SMS via Message API. Its status is available by `list` API

### Flow:

By send message, the `Transactional Outbox` will prepare a copy of desired message to handle its `sending` flow
by sharing the message across the `Kafka` topics and the `goroutine workers` support the process in background.
It has the `retry` and `dead-letter queue` to guarantee the sending and achieve the final state.

- The logs are handled in 3 levels of `Terminal Stdout`, `File` and `Logstash(ELK)`. 
- The `Mocked SMS Provider` traced by `Opentelemetry` and `Jaeger`
- The `Kafka` topics progress are available by `Kafka-UI`
- The APIs are available by `Swagger(OpenAPI)`

---

## Quick Start

To start service by docker-compose:
```shell
make run
```

Note: due to high resource, wait to have all services.(use `make up` for failed ones)

To Stop:
```shell
make down
```
---

The documented APIs are accessible at the link below:

```text
http://0.0.0.0:8080/public/swagger/index.html
```
Credentials
```text
username: admin
password: admin
```
---

The trace dashboard:

```text
http://0.0.0.0:16686/search
```

---

The queue dashboard:

```text
http://0.0.0.0:8082
```

---

Note: file logs are available in `./docker/sms-logs`


---