FROM alpine

WORKDIR /app

COPY --from=statistics_service_build:develop /app/cmd/statistics/main ./app

COPY clickhouse ./clickhouse
COPY .config.yaml .

CMD ["/app/app"]