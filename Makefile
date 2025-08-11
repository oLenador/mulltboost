# Variáveis reutilizáveis
DC := docker compose
COMPOSE_FILE := ./docker/docker-compose.dev.yaml

.PHONY: help dev-up dev-down dev-stop dev-restart dev-logs dev-ps dev-shell

help:
	@echo "Targets disponíveis:"
	@echo "  dev-up       - sobe o container (detached)"
	@echo "  dev-down     - encerra e remove container"
	@echo "  dev-stop     - apenas para container"
	@echo "  dev-restart  - rebuild e up"
	@echo "  dev-logs     - exibe logs em tempo real"
	@echo "  dev-ps       - lista containers"
	@echo "  dev-shell    - abre shell no container windows (se aplicável)"

dev-up:
	$(DC) -f $(COMPOSE_FILE) up 

dev-down:
	$(DC) -f $(COMPOSE_FILE) down

dev-stop:
	$(DC) -f $(COMPOSE_FILE) stop

dev-restart: down up

dev-logs:
	$(DC) -f $(COMPOSE_FILE) logs -f

dev-ps:
	$(DC) -f $(COMPOSE_FILE) ps

dev-shell:
	$(DC) -f $(COMPOSE_FILE) exec windows bash
