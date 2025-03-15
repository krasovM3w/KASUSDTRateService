# KAS_USDT Rate Service

Микросервис для получения курса KAS к USDT с биржи MEXC и сохранения данных в PostgreSQL

## Особенности

- gRPC API (GetRate, HealthCheck)
- Сохранение курса в PostgreSQL
- Мониторинг через Prometheus
- Трассировка через OpenTelemetry + Jaeger
- Конфигурация через флаги и переменные окружения
- Graceful shutdown
- Миграции базы данных
- Логирование через Zap

## Требования

- Go 1.22+
- Docker
- Protoc
- PostgreSQL 14+
- Prometheus (опционально)
- Jaeger (опционально)

## Быстрый старт

### Запуск с Docker Compose

```bash
# Скопируйте пример конфигурации
cp .env.example .env

# Запустите сервисы
docker-compose -f deployments/docker-compose.yaml up -d

# Примените миграции
make migrate-up

# Запустите приложение
make run

# Запустите тесты
make test
```

## Использование API

```bash
# Получить текущий курс
grpcurl -plaintext -d '{"base_currency": "KAS", "target_currency": "USDT"}' \
  localhost:50051 currensy.CurrencyService/GetRate
  
# Health Check
grpcurl -plaintext localhost:50051 currensy.CurrencyService/HealthCheck

# Нагрузочный тест
#cmd
for /L %i in (1,1,5000) do grpcurl -plaintext -d "{\"base_currency\":\"KAS\",\"target_currency\":\"USDT\"}" localhost:50051 currensy.CurrencyService.GetRate

#pShell
1..5000 | ForEach-Object { grpcurl -plaintext -d '{"base_currency":"KAS","target_currency":"USDT"}' localhost:50051 currensy.CurrencyService.GetRate }
```
## Приоритет конфигурации: 

Флаги > Переменные окружения > Значения по умолчанию

| Параметр           | Переменная окружения | По умолчанию                          | Описание                              |
|---------------------|----------------------|---------------------------------------|---------------------------------------|
| `--db-url`          | `DB_URL`             | `postgres://...`                      | URL подключения к PostgreSQL          |
| `--grpc-port`       | `GRPC_PORT`          | `50051`                               | Порт gRPC сервера                     |
| `--metrics-port`    | `METRICS_PORT`       | `9090`                                | Порт для Prometheus метрик            |
| `--jaeger-endpoint` | `JAEGER_ENDPOINT`    | `http://localhost:14268/api/traces`   | Адрес Jaeger Collector                |
| `--mexc-timeout`    | `MEXC_TIMEOUT`       | `5s`                                  | Таймаут запросов к MEXC API           |
| `--service-name`    | `SERVICE_NAME`       | `usdt-rate-service`                   | Имя сервиса для трейсинга             |#   K A S U S D T R a t e S e r v i c e 
 
 
