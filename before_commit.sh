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

echo "gom test"
gom test plaintodo/*.go
gom test plaintodo/query/*.go
gom test plaintodo/config/*.go
gom test plaintodo/util/*.go
gom test plaintodo/executor/*.go
gom test plaintodo/command/*.go
gom test plaintodo/ls/*.go
gom test plaintodo/output/*.go
