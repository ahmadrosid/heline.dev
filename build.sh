#!/bin/bash
cd ui
yarn install
yarn build
cd ..
go build
./heline server start