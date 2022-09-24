# API Gateway

## Microservice API Gateway for routing between different services by url slug

## Prerequisites

- Go version 1.19
-

## Usage

Define your services inside `services.yaml`

Example;

```
order-api:
  host: "http://order-api.order-api.svc.cluster.local"
  name: "order-api"
```

## Testing

Run unit tests

```
make unit-test
```

## Linting

```
make lint
```
