FROM alpine:3.17
WORKDIR /listener
COPY listenerApp .
CMD ["./listenerApp"]