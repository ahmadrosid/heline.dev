#!/bin/bash

array=(
  "tokio-rs/tokio"
  "rails/rails"
  "laravel/framework"
  "laravel/laravel"
  "hashicorp/vault"
  "hashicorp/nomad"
  "hashicorp/vagrant"
  "hashicorp/terraform"
  "hashicorp/consul"
  "navidrome/navidrome"
  "bitwarden/web"
  "hairyhenderson/gomplate"
  "gookit/goutil"
  "sharkdp/bat"
  "TheAlgorithms/Go"
  "TheAlgorithms/Javascript"
  "TheAlgorithms/C-Plus-Plus"
  "TheAlgorithms/Python"
  "TheAlgorithms/Java"
)
for i in "${array[@]}"
do
	./heline scrape github $i;
done
