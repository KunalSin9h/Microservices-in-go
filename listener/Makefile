start: 
	@echo "Running Listener microservices"
	go run main.go

rabbitmq-up:
	@echo "Starting rabbitmq server"
	docker run --name rabbitmq -d --hostname localnode --rm \
		-p 5672:5672 -p 15672:15672 \
		rabbitmq:3.11-management-alpine
rabbitmq-down:
	@echo "Stoping rabbitmq server"
	docker kill rabbitmq
