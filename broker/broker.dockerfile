# [BUILD-STAGE] Base Go-Alpine Image used as Builder
FROM golang:1.19.5-alpine3.17 as builder
WORKDIR /broker
COPY . /broker/
# CGO_ENABLED=0 means we are not using any C libs.
RUN CGO_ENABLED=0 go build -o brokerApp ./cmd/api/*.go
RUN chmod +x /broker/brokerApp

# [DEPLOY-STAGE]
FROM alpine:3.17
WORKDIR /broker
COPY --from=builder /broker/brokerApp .
CMD ["./brokerApp"]
