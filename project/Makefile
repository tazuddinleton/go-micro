FRONT_END_BIN=frontApp
BROKER_BIN=brokerService
AUTH_BIN=authService
LOGGER_BIN=loggerService
MAIL_BINARY=mailerApp

## up: starts all containers in the background without forcing build
up:	
	@echo "Starting docker containers..."
	docker-compose up -d
	@echo "Docker container started."

## up_build: stops docker-compose, builds all projects and starts docker compose
up_build:	down build_broker build_auth build_logger build_mail
	@echo "Re building docker images (if necessary)..."
	docker-compose up --build -d
	@echo "Started."

down:	
	@echo "Stopping docker compose..."
	docker-compose down
	@echo "Done."

build_broker:	
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ${BROKER_BIN} ./cmd/api
	@echo "Done."

build_auth:
	@echo "Building auth service binary..."
	cd ../auth-service && env GOOS=linux CGO_ENABLED=0 go build -o ${AUTH_BIN} ./cmd/api
	@echo "Done."

build_logger:
	@echo "Building logger service binary..."
	cd ../logger-service && env GOOS=linux CGO_ENABLED=0 go build -o ${LOGGER_BIN} ./cmd/api
	@echo "Done."

## build_mail: builds the mail binary as a linux executable
build_mail:
	@echo "Building mail binary..."
	cd ../mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ${MAIL_BINARY} ./cmd/api
	@echo "Done!"

build_front:
	@echo "Building frontend binary..."
	cd ../front-end && go build -o ${FRONT_END_BIN} ./cmd/web	
	@echo "Done."

start:
	@echo "Starting frontend..."    
	cd ../front-end && ./${FRONT_END_BIN} &
	@echo "Done."

stop: 
	@echo "Stopping frontend..."
	@-pkill -SIGTERM -f "./${FRONT_END_BIN}"
	@echo "Done."
