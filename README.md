# Web crawler

## What

A web crawling system

- Simple API to POST URLs for crawling and GET the results of a crawl (Go)
- URLs for crawling are queued up (RabbitMQ)
- Queue consumer written in Go to scrape web pages and queue up scraped URLs (Go)
- Browser instance to render web pages ([go-rod](https://github.com/go-rod/rod))
- Caching layer to reduce checks for recently crawled URLs (Redis)
- Database to persist crawl events (Postgres)
- Development environment (Docker/Docker Compose)

## Why

1. Learn Go and Kubernetes
1. General exploration

## Getting started

### Requirements

- docker v25
- docker-compose v2.24
- go 1.22

### Commands

| Description                                       | Command                         |
| ------------------------------------------------- | ------------------------------- |
| Run detached                                      | `docker-compose up -d`          |
| Build + run                                       | `docker-compose up --build`     |
| Down                                              | `docker-compose down`           |
| Down and also remove volumes (start fresh on dbs) | `docker-compose down --volumes` |

## Usage

| Description            | Command                                                                                                              |
| ---------------------- | -------------------------------------------------------------------------------------------------------------------- |
| Add a URL for crawling | `curl -d '{"url":"https://leevincenth.com"}' -H "Content-Type: application/json" -X POST http://localhost:80/crawls` |
| Postgres               | docker exec -it crawler-db psql -U <PGUSER> -d <PGDATABASE>                                                          |
| Logs                   | docker logs <crawler-worker,crawler-api, etc>                                                                        |
| View Queue             | Go to http://localhost:15672                                                                                         |

## Todo

### Core features

- [x] crawler
- [x] persistence
- [x] queuing
- [ ] tests

### Build

- [x] docker
- [x] docker compose

### Deploy

- [ ] CI/CD pipelines
- [ ] kubernetes

### Observability

- [ ] metrics
- [ ] alerting
- [ ] ELK
