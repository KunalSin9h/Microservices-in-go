FROM alpine:3.17
WORKDIR /frontend
COPY frontendApp .
CMD ["./frontendApp"]
