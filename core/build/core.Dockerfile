FROM alpine

WORKDIR /app

COPY --from=core_service_build:develop /app/cmd/core/main ./app

COPY postgresql ./postgresql
COPY .config.yaml .

CMD ["/app/app"]