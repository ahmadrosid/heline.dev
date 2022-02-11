#!/bin/bash
cd ui
yarn build
cd ..
go build
./heline server start