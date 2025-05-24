#!/bin/bash

# Source the environment file if it exists and we're not in Docker
if [ -f "scripts/env.sh" ]; then
  . scripts/env.sh
fi

# In Docker, these environment variables should be set by docker-compose
# If not set, use defaults (for local development)
SOLR_PORT=${SOLR_PORT:-8984}
SOLR_BASE_URL=${SOLR_BASE_URL:-"http://localhost:$SOLR_PORT"}

# Print Solr configuration for debugging
echo "Using Solr at: $SOLR_BASE_URL"

START="start"
STOP="stop"
PREPARE="prepare"
CLEAN="clean"

if [ "$1" == $START ]; then
  # Check if we're in Docker
  if [ -n "$DOCKER_ENV" ]; then
    echo "In Docker environment, Solr should be started by Docker Compose"
    # Check if Solr is accessible
    if curl -s -f "$SOLR_BASE_URL/solr/admin/info/system" > /dev/null; then
      echo "Solr is running at $SOLR_BASE_URL"
    else
      echo "WARNING: Cannot connect to Solr at $SOLR_BASE_URL"
    fi
  else
    # Local environment
    # Don't start solr if it already run.
    if [ -d "$SOLR_BUILD_FOLDER" ] && [ -f "$SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/bin/solr" ]; then
      bash $SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/bin/solr status
      if [ "$?" == "0" ]; then
        echo "${blue}Solr server already started.${normal}"
      else
        echo "${green}Starting solr server...${normal}"
        sudo bash $SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/bin/solr start -p $SOLR_PORT -force
      fi
    else
      echo "Solr installation not found at $SOLR_BUILD_FOLDER/solr-$SOLR_VERSION"
    fi
  fi
fi

if [ "$1" == $STOP ]; then
  # Check if we're in Docker
  if [ -n "$DOCKER_ENV" ]; then
    echo "In Docker environment, Solr should be stopped by Docker Compose"
  else
    echo "${green}Stopping solr server...${normal}"
    if [ -d "$SOLR_BUILD_FOLDER" ] && [ -f "$SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/bin/solr" ]; then
      sudo bash $SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/bin/solr stop -force -all
    else
      echo "Solr installation not found at $SOLR_BUILD_FOLDER/solr-$SOLR_VERSION"
    fi
  fi
fi

if [ "$1" == $PREPARE ]; then
  # Copy initial solr config file for heline index if not exists.
  if [ -n "$DOCKER_ENV" ]; then
    echo "In Docker environment, Solr core is created by Docker Compose"
  elif ! test -d "$SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/server/solr/heline"; then
    # Use sudo only if not in Docker
    cp -r "$SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/server/solr/configsets/_default" "$SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/server/solr/heline" || \
    sudo cp -r "$SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/server/solr/configsets/_default" "$SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/server/solr/heline"
  fi

  # Create core
  curl --request GET \
    --url "$SOLR_BASE_URL/solr/admin/cores?action=CREATE&name=heline&instanceDir=heline&config=solrconfig.xml&dataDir=data"

  # Create code_syntax field type for better handling of code patterns
  curl --request POST \
    --url "$SOLR_BASE_URL/solr/heline/schema" \
    --header 'Accept: application/json' \
    --header 'Content-type: application/json' \
    --data '{
    "add-field-type": {
      "name": "code_syntax",
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
          "class": "solr.ClassicTokenizerFactory"
        },
        "filters": [
          {
            "class": "solr.LowerCaseFilterFactory"
          },
          {
            "class": "solr.ShingleFilterFactory",
            "minShingleSize": "2",
            "maxShingleSize": "5",
            "outputUnigrams": "true"
          },
          {
            "class": "solr.RemoveDuplicatesTokenFilterFactory"
          }
        ]
      },
      "query": {
        "charFilters": [
          {
            "class": "solr.HTMLStripCharFilterFactory"
          }
        ],
        "tokenizer": {
          "class": "solr.ClassicTokenizerFactory"
        },
        "filters": [
          {
            "class": "solr.LowerCaseFilterFactory"
          },
          {
            "class": "solr.SynonymGraphFilterFactory",
            "expand": "true",
            "ignoreCase": "true",
            "synonyms": "synonyms.txt"
          }
        ]
      }
    }
  }'
  
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
          },
          {
            "class": "solr.PatternReplaceCharFilterFactory",
            "pattern": "([\\p{Punct}&&[^_]])",
            "replacement": " $1 "
          }
        ],
        "tokenizer":{
          "class": "solr.WhitespaceTokenizerFactory",
          "rule": "java"
        },
        "filters": [
          {
            "class":"solr.WordDelimiterFilterFactory",
            "generateWordParts": "1",
            "generateNumberParts": "1",
            "catenateWords": "1",
            "catenateNumbers": "1",
            "catenateAll": "0",
            "splitOnCaseChange": "1",
            "preserveOriginal": "1"
          },
          {
            "class": "solr.LowerCaseFilterFactory"
          },
          {
            "class":"solr.ASCIIFoldingFilterFactory"
          },
          {
            "class": "solr.StopFilterFactory",
            "ignoreCase": "true",
            "words": "stopwords.txt"
          }
        ]
      },
      "query": {
        "charFilters": [
          {
            "class": "solr.HTMLStripCharFilterFactory"
          },
          {
            "class": "solr.PatternReplaceCharFilterFactory",
            "pattern": "([\\p{Punct}&&[^_]])",
            "replacement": " $1 "
          }
        ],
        "tokenizer": {
          "class": "solr.WhitespaceTokenizerFactory",
          "rule": "java"
        },
        "filters": [
          {
            "class":"solr.WordDelimiterFilterFactory",
            "generateWordParts": "1",
            "generateNumberParts": "1",
            "catenateWords": "1",
            "catenateNumbers": "1",
            "catenateAll": "0",
            "splitOnCaseChange": "1",
            "preserveOriginal": "1"
          },
          {
            "class": "solr.LowerCaseFilterFactory"
          },
          {
            "class":"solr.ASCIIFoldingFilterFactory"
          },
          {
            "class": "solr.StopFilterFactory",
            "ignoreCase": "true",
            "words": "stopwords.txt"
          }
        ]
      }
    }
  }'
  
  # Create text_ngram field type for partial matching
  curl --request POST \
    --url "$SOLR_BASE_URL/solr/heline/schema" \
    --header 'Accept: application/json' \
    --header 'Content-type: application/json' \
    --data '{
    "add-field-type": {
      "name": "text_ngram",
      "class": "solr.TextField",
      "positionIncrementGap": "100",
      "analyzer": {
        "charFilters": [
          {
            "class": "solr.PatternReplaceCharFilterFactory",
            "pattern": "([\\p{Punct}&&[^_]])",
            "replacement": " $1 "
          }
        ],
        "tokenizer":{
          "class": "solr.NGramTokenizerFactory",
          "minGramSize": "2",
          "maxGramSize": "15"
        },
        "filters": [
          {
            "class": "solr.LowerCaseFilterFactory"
          }
        ]
      },
      "query": {
        "charFilters": [
          {
            "class": "solr.PatternReplaceCharFilterFactory",
            "pattern": "([\\p{Punct}&&[^_]])",
            "replacement": " $1 "
          }
        ],
        "tokenizer": {
          "class": "solr.StandardTokenizerFactory"
        },
        "filters": [
          {
            "class": "solr.LowerCaseFilterFactory"
          }
        ]
      }
    }
  }'

  # Create code_pattern field type for exact code syntax matching
  curl --request POST \
    --url "$SOLR_BASE_URL/solr/heline/schema" \
    --header 'Accept: application/json' \
    --header 'Content-type: application/json' \
    --data '{
    "add-field-type": {
      "name": "code_pattern",
      "class": "solr.TextField",
      "positionIncrementGap": "100",
      "analyzer": {
        "tokenizer":{
          "class": "solr.KeywordTokenizerFactory"
        },
        "filters": [
          {
            "class": "solr.LowerCaseFilterFactory"
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
    },
    "add-field": {
      "name": "code_content",
      "type": "code_syntax",
      "multiValued": true,
      "stored": true,
      "indexed": true
    },
    "add-field": {
      "name": "identifier_ngram",
      "type": "text_ngram",
      "stored": true,
      "indexed": true
    },
    "add-field": {
      "name": "code_syntax",
      "type": "code_pattern",
      "stored": true,
      "indexed": true
    }
  }'
fi

if [ "$1" == $CLEAN ]; then
  # Delete core first
  curl --request GET \
    --url "$SOLR_BASE_URL/solr/admin/cores?action=UNLOAD&core=heline&deleteIndex=true&deleteDataDir=true&deleteInstanceDir=true"

  # Delete solr folder
  if [ -n "$DOCKER_ENV" ]; then
    echo "In Docker environment, Solr data is managed by Docker volumes"
  else
    # Try without sudo first, then with sudo if needed
    rm -rf "$SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/server/solr/heline" || \
    sudo rm -rf "$SOLR_BUILD_FOLDER/solr-$SOLR_VERSION/server/solr/heline"
  fi
fi
