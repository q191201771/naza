#!/usr/bin/env bash

set -e

gofiles=$(git diff --name-only --diff-filter=ACM | grep '.go$')
unformatted=$(gofmt -l $gofiles)

if [ ! -z "$unformatted" ]; then
  echo >&2 "Go files should be formatted with gofmt. Please run:"
  for fn in $unformatted; do
      echo >&2 "  gofmt -w $PWD/$fn"
  done
  #exit 1
else
  echo >&2 "Go files be formatted."
fi

echo "" > coverage.txt

for d in $(go list ./... | grep -v vendor | grep nezha/pkg); do
    go test -race -coverprofile=profile.out -covermode=atomic $d
    if [ -f profile.out ]; then
        cat profile.out >> coverage.txt
        rm profile.out
    fi
done

# go test -race -coverprofile=profile.out -covermode=atomic && go tool cover -html=profile.out -o coverage.html && open coverage.html
# go test -test.bench=".*"
# go test -bench=. -benchtime=10s
