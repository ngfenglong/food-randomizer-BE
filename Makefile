build:
	@echo Building...
	@go build -o dist/api ./cmd/api
	@echo Backend built!

start: build
	@echo Starting the Backend ...
	./dist/api
	@echo Backend running!

stop:
	@echo Stopping the Backend...
	@pkill -f "./dist/api"
	@echo Stopped Backend