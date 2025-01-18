# cd .\linux\project
# make -f makefile_my_win.mak build_auth
# make -f makefile_my_win.mak build_broker
# make -f makefile_my_win.mak up_build
# make -f makefile_my_win.mak up
# curl http://localhost:8080/
# make -f makefile_my_win.mak down

# make -f makefile_my_win.mak build_front
# make -f makefile_my_win.mak start
# http://localhost:8080/
# make -f makefile_my_win.mak stop
#
#type makefile_my_win.mak | findstr ":"

SHELL=cmd.exe
FRONT_END_BINARY=feApp
BROKER_BINARY=brokerApp
AUTH_BINARY=authApp
TRACE_BINARY=traceApp
MAIL_BINARY=mailApp
LISTENER_BINARY=listenerApp
DOCKER_COMPOSE_PATH=..\project\docker-compose.yaml

up: ## starts all containers in the background without forcing build
	@echo Stopping docker images (if running...)
	docker-compose -f ${DOCKER_COMPOSE_PATH} down
	@echo Starting Docker images...
	docker-compose -f ${DOCKER_COMPOSE_PATH} up -d
	@echo Docker images started!

build: build_broker build_auth build_mail build_listener build_trace build_fe
	@echo "All Docker images built"

up_build: build_broker build_auth build_trace build_mail build_listener build_fe ## stops docker-compose (if running), builds all projects and starts docker compose
	@echo Stopping docker images (if running...)
	docker-compose -f ${DOCKER_COMPOSE_PATH} down
	@echo Building (when required) and starting docker images...
	docker-compose -f ${DOCKER_COMPOSE_PATH} up --build -d
	@echo Docker images built and started!

down: ## down stop docker compose
	@echo Stopping docker compose...
	docker-compose -f ${DOCKER_COMPOSE_PATH} down
	@echo Done!

build_broker: ##  builds the broker binary as a linux executable
	@echo Building broker binary...
	chdir ..\broker-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o .\app\${BROKER_BINARY} ./cmd/api && dir .\app\${BROKER_BINARY}
	@echo Done!

###################################################################################################

build_auth:
	@echo Building auth binary...
	chdir ..\authentication-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o .\app\${AUTH_BINARY} ./cmd/api && dir .\app\${AUTH_BINARY}
	@echo Done!

#################################################################################################

build_mail:
	@echo Building mail binary...
	chdir ..\mail-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o .\app\${MAIL_BINARY} ./cmd/api && dir .\app\${MAIL_BINARY}
	@echo Done!

#################################################################################################

build_listener:
	@echo Building listener binary...
	chdir ..\listener-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o .\app\${LISTENER_BINARY} . && dir .\app\${LISTENER_BINARY}
	@echo Done!

#################################################################################################

build_trace:
	@echo Building auth binary...
	chdir ..\trace-service && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o .\app\${TRACE_BINARY} ./cmd/api && dir .\app\${TRACE_BINARY}
	@echo Done!

#################################################################################################
## -v print the names of packages as they are compiled. -o output -a force rebuilding of packages that are already up-to-date.
## builds front-end service
build_fe:
	@echo "Building fe binary "
	cd ..\fe && set GOOS=linux&& set GOARCH=amd64&& set CGO_ENABLED=0 && go build -o ./app/${FRONT_END_BINARY} ./cmd/web && dir .\app\${FRONT_END_BINARY}
	@echo "Done!"

#################################################################################################

build_front:  # builds the frone end binary
	@echo Building front end binary...
	chdir ..\fe && set CGO_ENABLED=0&& set GOOS=windows&& go build -v -a -o ${FRONT_END_BINARY} ./cmd/web && dir ${FRONT_END_BINARY}
	@echo Done!

start: build_front  # starts the front end
	@echo Starting front end
	chdir ..\fe && start /B ${FRONT_END_BINARY} &

stop: # stop the front end
	@echo Stopping front end...
	@taskkill /IM "${FRONT_END_BINARY}" /F
	@echo "Stopped front end!"
