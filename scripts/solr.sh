#!/bin/bash
. scripts/env.sh

START="start"
STOP="stop"
PREPARE="prepare"
CLEAN="clean"

if [ "$1" == $START ]; then
  # Don't start solr if it already run.
  bash $SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/bin/solr status
  if [ "$?" == "0" ]; then
    echo "${blue}Solr server already started.${normal}"
  else
    echo "${green}Starting solr server...${normal}"
    sudo bash $SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/bin/solr start -p $SOLR_PORT -force
  fi
fi

if [ "$1" == $STOP ]; then
  echo "${green}Stoping solr server...${normal}"
  sudo bash $SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/bin/solr stop -force -all
fi

if [ "$1" == $PREPARE ]; then
  # Copy initial solr config file for heline index if not exists.
  if ! test -d "$SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/server/solr/heline"; then
    sudo cp -r "$SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/server/solr/configsets/_default" "$SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/server/solr/heline"
  fi

  # Copy initial solr config file for docset index if not exists.
  if ! test -d "$SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/server/solr/docset"; then
    sudo cp -r "$SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/server/solr/configsets/_default" "$SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/server/solr/docset"
  fi

  # Create core
  curl --request GET \
    --url "$SOLR_BASE_URL/solr/admin/cores?action=CREATE&name=heline&instanceDir=heline&config=solrconfig.xml&dataDir=data"

  # Create core
  curl --request GET \
    --url "$SOLR_BASE_URL/solr/admin/cores?action=CREATE&name=docset&instanceDir=docset&config=solrconfig.xml&dataDir=data"

  # Create text_html field schema
  curl --request POST \
    --url "$SOLR_BASE_URL/solr/heline/schema" \
    --header 'Accept: application/json' \
    --header 'Content-type: application/json' \
    --data '{
    "add-field-type": {
      "name": "text_html",
      "class": "solr.TextField",
      "positionIncrementGap": "100",
      "autoGeneratePhraseQueries": "true",
      "analyzer": {
        "charFilters": [
          {
            "class": "solr.HTMLStripCharFilterFactory"
          }
        ],
        "tokenizer":{
          "class": "solr.WhitespaceTokenizerFactory",
          "rule": "java"
        },
        "tokenizer":{
          "class": "solr.NGramTokenizerFactory"
        },
        "filters": [
          {
            "class":"solr.WordDelimiterFilterFactory"
          },
          {
            "class": "solr.LowerCaseFilterFactory"
          },
          {
            "class":"solr.ASCIIFoldingFilterFactory"
          }
        ]
      },
      "query": {
        "tokenizer": {
          "class": "solr.WhitespaceTokenizerFactory",
          "rule": "java"
        },
        "filters": [
          {
            "class":"solr.WordDelimiterFilterFactory"
          },
          {
            "class": "solr.LowerCaseFilterFactory"
          },
          {
            "class":"solr.ASCIIFoldingFilterFactory"
          }
        ]
      }
    }
  }'

  # Create index schema
  curl --request POST \
    --url "$SOLR_BASE_URL/solr/heline/schema" \
    --header 'Content-type: application/json' \
    --data '{
    "add-field": {
      "name": "branch",
      "type": "string",
      "stored": true
    },
    "add-field": {
      "name": "path",
      "type": "string",
      "stored": true
    },
    "add-field": {
      "name": "file_id",
      "type": "string",
      "stored": true
    },
    "add-field": {
      "name": "owner_id",
      "type": "string",
      "stored": true
    },
    "add-field": {
      "name": "lang",
      "type": "string",
      "stored": true
    },
    "add-field": {
      "name": "repo",
      "type": "string",
      "stored": true
    },
    "add-field": {
      "name": "content",
      "type": "text_html",
      "multiValued": true,
      "stored": true,
      "indexed": true
    }
  }'
fi

if [ "$1" == $CLEAN ]; then

  # 1. Delete all index
  curl --request GET \
    --url "$SOLR_BASE_URL/solr/heline/update?commit=true" \
    --header 'Content-Type: application/json' \
    --data '{
    "delete": {
      "query": "*:*"
    }
  }'

  # 2. Delete scheme
  curl --request POST \
    --url "$SOLR_BASE_URL/solr/heline/schema" \
    --header 'Content-type: application/json' \
    --data '{
    "delete-field": {
      "name": "content"
    },
    "delete-field": {
      "name": "branch"
    },
    "delete-field": {
      "name": "path"
    },
    "delete-field": {
      "name": "file_id"
    },
    "delete-field": {
      "name": "owner_id"
    },
    "delete-field": {
      "name": "lang"
    },
    "delete-field": {
      "name": "repo"
    }
  }'

  # 3. Delete field text_html
  curl --request POST \
    --url "$SOLR_BASE_URL/solr/heline/schema" \
    --header 'Content-type: application/json' \
    --data '{
    "delete-field-type": {
      "name": "text_html"
    }
  }'

  # Delete core
  curl --request GET \
    --url "$SOLR_BASE_URL/solr/admin/cores?action=UNLOAD&core=heline"


  # Delete core
  curl --request GET \
    --url "$SOLR_BASE_URL/solr/admin/cores?action=UNLOAD&core=docset"

  # Delete solr folder
  sudo rm -rf "$SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/server/solr/heline"
  sudo rm -rf "$SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/server/solr/docset"
fi
