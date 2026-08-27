FROM golang:1.23 AS builder
ENV GOPROXY=off GOSUMDB=off
WORKDIR /src
COPY go.mod go.sum ./
COPY vendor ./vendor
COPY . .
RUN go build -mod=vendor -o /out/windturbinectl ./cmd/wtc

FROM golang:1.23
ENV GOPROXY=off GOSUMDB=off
WORKDIR /app
COPY --from=builder /src ./
COPY --from=builder /out/windturbinectl /app/bin/windturbinectl
RUN mkdir -p /app/data
EXPOSE 8080
CMD ["/app/bin/windturbinectl"]
