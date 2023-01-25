FROM alpine:3.17
WORKDIR /auth
COPY authApp .
CMD ["./authApp"]