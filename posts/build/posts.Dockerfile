FROM alpine

WORKDIR /app

COPY --from=posts_service_build:develop /app/cmd/posts/main ./app

COPY postgresql ./postgresql
COPY .config.yaml .

CMD ["/app/app"]