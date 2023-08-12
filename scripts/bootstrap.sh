#!/bin/bash

set -e

. scripts/env.sh

if ! test -d "$SOLR_BUILD_FOLDER"; then
  mkdir $BUILD_FOLDER
fi

# Setup solr
if ! test -d "$SOLR_BUILD_FOLDER"; then
  mkdir $SOLR_BUILD_FOLDER
  echo "Folder $SOLR_BUILD_FOLDER does not exist, downloading..."
  wget -O "$SOLR_BUILD_FOLDER/$SOLR_VERSION.tgz" $SOLR_DOWNLOAD_URL
  tar zxvf "$SOLR_BUILD_FOLDER/$SOLR_VERSION.tgz"
  mv ./solr-$SOLR_VERSION $SOLR_BUILD_FOLDER
fi

# Install java if not exists
if !(command -v java); then
  sudo add-apt-repository ppa:openjdk-r/ppa
  sudo apt-get update
  sudo apt install openjdk-11-jdk -y
fi

# Make sure to have go path
export PATH="$PATH:/usr/local/go/bin"

# Install go if not exists
if !(command -v go); then
  echo "Go binary not found, downloading..."
  curl $GO_DOWNLOAD_URL -o "$GO_BUILD_FOLDER/$GO_VERSION.tar.gz"
  tar zxvf "$GO_BUILD_FOLDER/$GO_VERSION.tar.gz"
  sudo mv ./go /usr/local/go
  echo 'export PATH="$PATH:/usr/local/go/bin"' >> ~/.bashrc
  source ~/.bashrc
fi

# Install nodejs
if !(command -v node); then
  curl -sL https://deb.nodesource.com/setup_16.x | sudo -E bash -
  sudo apt install nodejs -y
  npm install --global pnpm
fi

# Install rust
if !(command -v cargo); then
  curl --proto '=https' --tlsv1.2 -sSf https://sh.rustup.rs | sh
  cd $BUILD_FOLDER
  git clone https://github.com/ahmadrosid/heline-indexer.git
fi
