FROM golang:1.22-alpine AS build
WORKDIR /app
COPY . .
RUN go build -o portal ./cmd/portal

FROM alpine:3.20
WORKDIR /srv
COPY --from=build /app/portal /usr/local/bin/portal
COPY assets /srv/assets
COPY templates /srv/templates
EXPOSE 8080
ENTRYPOINT ["portal"]
