#!/bin/bash
cd ui
if [ ! -d "node_modules" ]; then
    pnpm install
fi

cd ..
go build
./heline server start
