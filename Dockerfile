FROM golang:1.20-alpine AS builder
WORKDIR /app

COPY . ./

RUN go build -o /aggregator ./cmd/aggregator/main.go


FROM scratch

WORKDIR /

COPY --from=builder /aggregator /aggregator

CMD ["/aggregator"]

