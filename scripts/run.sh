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
    cd $UI_FOLDER && pnpm dev
  else
    cd $UI_FOLDER && pnpm && pnpm dev
  fi
}

if [ "$1" == "ui" ]; then
  start_ui
fi

if [ "$1" == "solr start" ]; then
  start_ui
fi

if [ "$1" == "stop" ]; then
  lsof -i tcp:8984 | grep LISTEN | awk '{print $2}' | xargs kill -9
  lsof -i tcp:8000 | grep LISTEN | awk '{print $2}' | xargs kill -9
fi

if [ "$1" == "reset" ]; then
  set -ex
  . scripts/solr.sh clean
  . scripts/solr.sh prepare
fi

if [ "$1" == "build" ]; then
  cd ui && pnpm build
  cd ../ && go build
fi

if [ "$1" == "production" ]; then
  . scripts/solr.sh start
  ./heline &>/dev/null & disown;
fi

if [ "$1" == "gen" ]; then
  cd ui && pnpm build
  cd ..
  git add .
  git commit -m 'chore: regenerate static file'
  git push origin HEAD
fi
