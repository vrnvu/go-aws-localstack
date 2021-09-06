.PHONY = help stop restart ps

.DEFAULT_GOAL = help

help:
	@echo "---------------HELP-----------------"
	@echo "stop: docker-compose stop"
	@echo "restart: docker-compose restart"
	@echo "ps: docker-compose ps"
	@echo "------------------------------------"

stop:
	docker-compose stop

restart:
	docker-compose restart

ps:
	docker-compose ps
