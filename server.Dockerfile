FROM golang:1.20-alpine AS builder
WORKDIR /app

COPY . ./

RUN go build -o /aggregator ./cmd/aggregator/main.go


FROM alpine

WORKDIR /

COPY --from=builder /aggregator /aggregator
COPY --from=builder /app/conf.yml /conf.yml

EXPOSE 8080

CMD ["/aggregator"]

