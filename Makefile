.PHONY = help run stop restart ps

.DEFAULT_GOAL = help

help:
	@echo "---------------HELP-----------------"
	@echo "run: execute the main.go"
	@echo "stop: docker-compose stop"
	@echo "restart: docker-compose restart"
	@echo "ps: docker-compose ps"
	@echo "------------------------------------"

run:
	go run main.go

stop:
	docker-compose stop

restart:
	docker-compose restart

ps:
	docker-compose ps
