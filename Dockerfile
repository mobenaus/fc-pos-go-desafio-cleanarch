FROM golang:latest AS builder
WORKDIR /app
COPY . .
RUN cd cmd/ordersystem && \
    GOOS=linux CGO_ENABLED=0 go build -o server main.go wire_gen.go

FROM scratch
COPY --from=builder /app/cmd/ordersystem/server .
COPY --from=builder /app/.env .
CMD ["./server"]