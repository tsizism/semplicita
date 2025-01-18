# make stop

FRONT_END_BINARY=feApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
TRACE_BINARY=traceApp
MAIL_BINARY=mailApp
LISTENER_BINARY=listenerApp

## up: starts all containers in the background without forcing build -d, --detach
up:
	@echo "Starting Docker images..."
	docker compose up -d
	@echo "Docker images started!"

## makes go execs
build: build_broker build_auth build_mail build_listener build_trace build_fe
	@echo "All Docker images built"

## builds docker images
image: 
	@echo "Building all Docker images"
	docker build --no-cache -f ../authentication-service/dockerfile -t tsizism/authentication-service:1.0 ../authentication-service
	docker build --no-cache -f ../broker-service/dockerfile -t tsizism/broker-service:1.0 ../broker-service
	docker build --no-cache -f ../listener-service/dockerfile -t tsizism/listener-service:1.0 ../listener-service
	docker build --no-cache -f ../mail-service/dockerfile -t tsizism/mail-service:1.0 ../mail-service
	docker build --no-cache -f ../trace-service/dockerfile -t tsizism/trace-service:1.0 ../trace-service
	docker build --no-cache -f ../fe/dockerfile -t tsizism/fe:1.0 ../fe
	@echo "Docker images built"


## up_build: stops docker-compose (if running), makes all projects and starts docker compose  , -d, --detach
up_build: build_broker build_auth build_mail build_listener build_trace build_fe
	@echo "Stopping docker images (if running...)"
	docker compose down
	@echo "Building (when required) and starting docker images..."
	docker compose up --build -d
	@echo "Docker images built and started!"

## down: stop docker compose
down:
	@echo "Stopping docker compose..."
	docker compose down
	@echo "Done!"

## -v print the names of packages as they are compiled
## -a force rebuilding of packages that are already up-to-date.
## -o output
## build_broker: builds the broker binary as a linux executable
build_broker:
	@echo "Building broker binary..."
	cd ../broker-service && env GOOS=linux CGO_ENABLED=0 go build -o ./app/${BROKER_BINARY} ./cmd/api
	@echo "Done!"

#################################################################################################

build_mail:
	@echo Building mail binary...
	cd ../mail-service && env GOOS=linux CGO_ENABLED=0 go build -o ./app/${MAIL_BINARY} ./cmd/api && ls -la ./app/${MAIL_BINARY}
	@echo Done!

#################################################################################################

build_auth:
	@echo Building auth binary...
	cd ../authentication-service && env GOOS=linux CGO_ENABLED=0 go build -o ./app/${AUTH_BINARY} ./cmd/api && ls -la ./app/${AUTH_BINARY}
	@echo Done!

#################################################################################################

build_listener:
	@echo Building listener binary...
	cd ../listener-service && env GOOS=linux CGO_ENABLED=0 go build -o ./app/${LISTENER_BINARY} . && ls -la ./app/${LISTENER_BINARY}
	@echo Done!

#################################################################################################

build_trace:
	@echo Building auth binary...
	cd ../trace-service && env GOOS=linux CGO_ENABLED=0  go build -o ./app/${TRACE_BINARY} ./cmd/api && ls -la ./app/${TRACE_BINARY}
	@echo Done!

#################################################################################################
## -v print the names of packages as they are compiled. -o output -a force rebuilding of packages that are already up-to-date.
## builds front-end service
build_fe:
	@echo "Building fe binary"
	cd ../fe && env GOOS=linux CGO_ENABLED=0 go build -a -o ./app/${FRONT_END_BINARY} ./cmd/web && ls -la ./app/${FRONT_END_BINARY}
	@echo "Done!"

#################################################################################################
## build_front: builds local front end binary
build_front:
	@echo "Building local front end binary..."
	cd ../fe && env CGO_ENABLED=0 && go build -v -a -o ${FRONT_END_BINARY} ./cmd/web && ls -la ${FRONT_END_BINARY}
	@echo "Done!"

## start: starts local front end
start: build_front
	@echo "Starting local front end"
	cd ../fe && ./${FRONT_END_BINARY} -port 8888 &


## stop: stop local front end
stop:
	@echo "Stopping local front end..."
	@-pkill -SIGTERM -f ${FRONT_END_BINARY}
	@echo "Stopped front end!"