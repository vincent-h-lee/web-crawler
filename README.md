# Web crawler

## What

A web crawling system

- Simple API to POST URLs for crawling and GET the results of a crawl
- Read URLs for crawling from a Queue (RabbitMQ)
- Caching layer to reduce checks for recently crawled URLs (Redis)
- Browser instance to load URLs ([go-rod](https://github.com/go-rod/rod))
- Database to persist crawl events (Postgres)
- Docker Compose for development

## Why

1. Learn Go and Kubernetes
1. General exploration

## Getting started

### Requirements

- docker v25
- docker-compose v2.24
- go 1.22

### Commands

| --- | --- |
| Description | Command |
| --- | --- |
| Run detached | docker-compose up -d |
| Build + run | docker-compose up --build |
| Down | docker-compose down |
| Down and also remove volumes (start fresh on dbs) | docker-compose down --volumes |
| Postgres | docker exec -it crawler-db psql -U <PGUSER> -d <PGDATABASE> |

## Todo

### Core features

- [x] crawler
- [x] persistence
- [x] queuing
- [] tests

### Build

- [x] docker
- [x] docker compose

### Deploy

- [] CI/CD pipelines
- [] kubernetes

### Observability

- [] metrics
- [] alerting
- [] ELK
