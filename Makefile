.PHONY: ui 

install:
	bash build.sh

ui:
	bash scripts/run.sh ui

gen:
	bash scripts/run.sh gen

go:
	go build && ./heline server start

reset:
	bash scripts/solr.sh clean && bash scripts/solr.sh prepare

prod:
	bash scripts/run.sh production

stop:
	bash scripts/run.sh stop
