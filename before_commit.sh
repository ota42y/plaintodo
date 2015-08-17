#!/bin/bash

echo "go fmt"
find . -name "*.go" -exec go fmt {} \;

echo "golint"
golint plaintodo/task
golint plaintodo/query
golint plaintodo/config
golint plaintodo/util
golint plaintodo/executor
golint plaintodo/command
golint plaintodo/ls
golint plaintodo/output
# golint palintodo
