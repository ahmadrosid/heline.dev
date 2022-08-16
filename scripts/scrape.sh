#!/bin/bash
array=(
  "gitlab.json"
  "sh.json"
)

for i in "${array[@]}"
do
	BASE_URL=http://localhost:8984 hli -- $i;
done
