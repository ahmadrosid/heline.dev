#!/bin/bash
cd ui
if [ ! -d "node_modules" ]; then
    pnpm install
fi

pnpm build
cd ..
go build
./heline server start
