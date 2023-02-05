FROM alpine:3.17
WORKDIR /logger
COPY loggerApp .
CMD ["./loggerApp"]