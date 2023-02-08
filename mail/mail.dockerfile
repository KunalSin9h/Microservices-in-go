FROM alpine:3.17
WORKDIR /mail
COPY mailApp .
COPY ./cmd/api/templates/ ./templates/
CMD ["./mailApp"]