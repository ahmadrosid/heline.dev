#!/bin/bash

function start_solr() {
. scripts/solr.sh start
}

function start_server() {
  cmd="./heline server start"
  $cmd &
}

function start_ui() {
  # Start UI server
  UI_FOLDER="$(pwd)/ui"
  if test -d "$UI_FOLDER/node_modules"; then
    cd $UI_FOLDER && yarn dev
  else
    cd $UI_FOLDER && yarn && yarn dev
  fi
}

if [ "$1" == "ui" ]; then
  start_ui
fi

if [ "$1" == "stop" ]; then
  lsof -i tcp:8984 | grep LISTEN | awk '{print $2}' | xargs kill -9
  lsof -i tcp:8000 | grep LISTEN | awk '{print $2}' | xargs kill -9
fi

if [ "$1" == "reset" ]; then
  set e
  . scripts/solr.sh clean
  . scripts/solr.sh prepare
  . scripts/scrape.sh
fi

if [ "$1" == "build" ]; then
  cd ui && yarn build
  cd ../ && go build
fi

if [ "$1" == "production" ]; then
  . scripts/solr.sh start
  ./heline server start &>/dev/null & disown;
fi

if [ "$1" == "scrape" ]; then
  ./heline scrape github $2 &>/dev/null & disown;
fi
