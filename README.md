# Web crawler

## Commands

| --- | --- |
| Description | Command |
| --- | --- |
| Compose Up | docker-compose up |
| Compose Up with build | docker-compose up --build |
| Compose Down | docker-compose down |
| Compose Down with volumes | docker-compose down --volumes |
| Psql | docker exec -it web-crawler-db-1 psql -U <PGUSER> -d <PGDATABASE> |

## Todo

### Core

- [x] crawler
- [] persistence
- [] queuing
- [] tests

### Infra

- [] docker
- [] kubernetes

### Vitals

- [] ELK/prometheus+grafana
