install:
	bash build.sh

go:
	go build && ./heline server start

reset:
	bash scripts/solr.sh clean && bash scripts/solr.sh prepare

prod:
	bash scripts/run.sh production
