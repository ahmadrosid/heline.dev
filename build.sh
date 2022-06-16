#!/bin/bash
cd ui
if [ ! -d "node_modules" ]; then
    yarn install
fi

yarn build
cd ..
go build
./heline server start
