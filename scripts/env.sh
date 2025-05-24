#!/bin/bash

BUILD_FOLDER="$(pwd)/_build"

SOLR_BUILD_FOLDER="$BUILD_FOLDER/solr"
SOLR_VERSION="8.11.4"
# SOLR_VERSION="7.7.3"
SOLR_DOWNLOAD_URL="https://dlcdn.apache.org/lucene/solr/$SOLR_VERSION/solr-$SOLR_VERSION.tgz"
# SOLR_DOWNLOAD_URL="https://dlcdn.apache.org/lucene/solr/8.11.1/solr-8.11.1-src.tgz"

GO_VERSION="go1.16.10"
GO_DOWNLOAD_URL="https://dl.google.com/go/$GO_VERSION.linux-amd64.tar.gz"
GO_BUILD_FOLDER="$BUILD_FOLDER/go"

# Use environment variables if set, otherwise use defaults
# This allows Docker to override these values
SOLR_PORT=${SOLR_PORT:-8984}

# In Docker, SOLR_BASE_URL will be set to http://solr:8983
# For local development, use localhost
SOLR_BASE_URL=${SOLR_BASE_URL:-"http://localhost:$SOLR_PORT"}

# Print Solr configuration for debugging
echo "Using Solr at: $SOLR_BASE_URL"

if test -t 1; then # if terminal
    ncolors=$(which tput > /dev/null && tput colors) # supports color
    if test -n "$ncolors" && test $ncolors -ge 8; then
        termcols=$(tput cols)
        bold="$(tput bold)"
        underline="$(tput smul)"
        standout="$(tput smso)"
        normal="$(tput sgr0)"
        black="$(tput setaf 0)"
        red="$(tput setaf 1)"
        green="$(tput setaf 2)"
        yellow="$(tput setaf 3)"
        blue="$(tput setaf 4)"
        magenta="$(tput setaf 5)"
        cyan="$(tput setaf 6)"
        white="$(tput setaf 7)"
    fi
fi
