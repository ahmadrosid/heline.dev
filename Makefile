install:
	bash build.sh

next:
	bash scripts/run.sh ui

go:
	go build && ./heline server start

reset:
	bash scripts/solr.sh clean && bash scripts/solr.sh prepare

prod:
	bash scripts/run.sh production

stop:
	bash scripts/run.sh stop
