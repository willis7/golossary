FROM alpine

RUN apk update && apk --no-cache add --update ca-certificates && update-ca-certificates
ADD app /app
ADD app_config.toml /app_config.toml
ENTRYPOINT ["/app"]
