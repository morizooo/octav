#!/bin/bash

set -e

if [ -z "$OCTAV_DB_NAME" ]; then
    OCTAV_DB_NAME=octav
fi

export GO_TAGS_OPT="-tags debug0"

if [ "$TRAVIS" == "true" ]; then
    echo " + Detected running under Travis CI"
    make glide
    make initdb
    make installdeps
fi

export OCTAV_TEST_DSN="root:@/$OCTAV_DB_NAME?parseTime=true"
export OCTAV_TRACE_DB=1
export OCTAV_DEBUG_FILE=/tmp/debug.out

exec make test
