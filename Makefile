install:
	bash build.sh

go:
	go build && ./heline server start

start:
	go build && ./heline server start && eval "cd ui && yarn dev"

reset:
	bash scripts/solr.sh clean && bash scripts/solr.sh prepare
