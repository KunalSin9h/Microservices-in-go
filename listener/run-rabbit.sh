docker run --name rabbitmq --hostname localnode \
    --rm -d -p 15672:15672 -p 5672:5672 \
    rabbitmq:3.11-management-alpine