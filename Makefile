.PHONY: build dev ui gen go reset prod stop index clean

# Binary name
BINARY_NAME=heline

# Output directory
TARGET_DIR=./

build:
	go build -o $(TARGET_DIR)$(BINARY_NAME)

dev:
	bash build.sh

ui:
	bash scripts/run.sh ui

gen:
	bash scripts/run.sh gen

go:
	go build && ./$(BINARY_NAME) server start

reset:
	bash scripts/solr.sh clean && bash scripts/solr.sh prepare

prod:
	bash scripts/run.sh production

stop:
	bash scripts/run.sh stop

index:
	bash scripts/indexing.sh

clean:
	rm -f $(TARGET_DIR)$(BINARY_NAME)
	rm -rf ui/dist
