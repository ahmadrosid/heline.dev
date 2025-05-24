.PHONY: build dev ui gen go reset prod stop stopall index clean docker-up docker-down docker-dev docker-index docker-reset

# Binary name
BINARY_NAME=heline

# Output directory
TARGET_DIR=./

build:
	go build -o $(TARGET_DIR)$(BINARY_NAME)

dev:
	bash build.sh & bash scripts/run.sh ui

ui:
	bash scripts/run.sh ui

gen:
	bash scripts/run.sh gen

reset:
	bash scripts/solr.sh clean && bash scripts/solr.sh prepare

prod:
	bash scripts/run.sh production

stop:
	bash scripts/run.sh stop

stopall:
	bash scripts/run.sh stop
	pkill -f "$(BINARY_NAME)" || true
	lsof -ti:3000 | xargs kill -9 || true
	@echo "All Heline.dev processes have been stopped"

index:
	bash scripts/indexing.sh

clean:
	rm -f $(TARGET_DIR)$(BINARY_NAME)
	rm -rf ui/dist

# Docker Compose targets
docker-up:
	docker-compose up -d

# Start the development environment with Docker
docker-dev:
	docker-compose up dev

# Run only the application service
docker-app:
	docker-compose up heline-app

# Run indexing in Docker
docker-index:
	docker-compose exec heline-app make index

# Reset Solr in Docker
docker-reset:
	docker-compose exec heline-app make reset

# Stop all Docker services
docker-down:
	docker-compose down
