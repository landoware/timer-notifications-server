FROM golang:1.25-alpine AS builder

WORKDIR /src
COPY go.mod go.sum ./
RUN go mod download
COPY . .
RUN CGO_ENABLED=0 go build -o /out/osrs-notifier-server ./...

FROM gcr.io/distroless/static-debian12:nonroot

COPY --from=builder /out/osrs-notifier-server /usr/local/bin/osrs-notifier-server

EXPOSE 8080
ENTRYPOINT ["/usr/local/bin/osrs-notifier-server"]
