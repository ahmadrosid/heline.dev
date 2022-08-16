#!/bin/bash
cd _build/heline-indexer

array=(
  "sh.json"
  # "tokio-rs/tokio"
  # "rails/rails"
  # "laravel/framework"
  # "laravel/laravel"
  # "hashicorp/vault"
  # "hashicorp/nomad"
  # "hashicorp/vagrant"
  # "hashicorp/terraform"
  # "hashicorp/consul"
  # "navidrome/navidrome"
  # "bitwarden/web"
  # "hairyhenderson/gomplate"
  # "gookit/goutil"
  # "sharkdp/bat"
  # "TheAlgorithms/Go"
  # "TheAlgorithms/Javascript"
  # "TheAlgorithms/C-Plus-Plus"
  # "TheAlgorithms/Python"
  # "TheAlgorithms/Java"
)
for i in "${array[@]}"
do
	BASE_URL=http://localhost:8984 hli -- $i;
done
