FRONT_END_BINARY=frontendApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
LOGGER_BINARY=loggerApp
MAIL_BINARY=mailApp
LISTENER_BINARY=listenerApp

## up: starts all containers in the background without forcing build
up:
	@echo "Starting Docker images..."
	docker-compose up -d
	@echo "Docker images started!"

## up_build: stops docker-compose (if running), builds all projects and starts docker compose
up_build: build_front build_broker build_auth build_logger build_mail build_listener
	@echo "Stopping docker images (if running...)"
	docker-compose down
	@echo "Building (when required) and starting docker images..."
	docker-compose up --build -d
	@echo "Docker images built and started!"

build_and_push_images_with_tag: build_front
	cd ../frontend && docker build -f frontend.dockerfile -t kunalsin9h/frontend-service:latest . && \
		docker push kunalsin9h/frontend-service:latest

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done!"

## build_front: builds the frone end binary
build_front:
	@echo "Building front end binary..."
	cd ../frontend && GOOS=linux CGO_ENABLED=0 go build -o ${FRONT_END_BINARY} ./cmd/web/*.go
	@echo "Done!"

## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker && GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BINARY} ./cmd/api/*.go
	@echo "Done!"

build_auth:
	@echo "Building auth binary..."
	cd ../auth && GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BINARY} ./cmd/api/*.go
	@echo "Done!"

build_logger:
	@echo "Building logger binary..."
	cd ../logger && GOOS=linux CGO_ENABLED=0 GIN_MODE=release go build -o ${LOGGER_BINARY} ./cmd/api/*.go
	@echo "Done!"

build_mail:
	@echo "Building mail binary..."
	cd ../mail && GOOS=linux CGO_ENABLED=0 GIN_MODE=release go build -o ${MAIL_BINARY} ./cmd/api/*.go
	@echo "Done!"

build_listener:
	@echo "Building listener binary..."
	cd ../listener && GOOS=linux CGO_ENABLED=0 go build -o ${LISTENER_BINARY} ./main.go
	@echo "Done!"

## start: starts the front end
start: build_front
	@echo "Starting front end"
	cd ../frontend && ./${FRONT_END_BINARY} &

## stop: stop the front end
stop:
	@echo "Stopping front end..."
	@-pkill -SIGTERM -f "./${FRONT_END_BINARY}"
	@echo "Stopped front end!"