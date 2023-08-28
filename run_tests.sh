#!/bin/bash

for arg in "$@"; do
  case $arg in
    -cover) coverage=true;;
    -gen) generate=true;;
  esac
done

if [[ "$generate" = true ]]; then
  ./generate_mocks.sh
fi

if [[ "$coverage" = true ]]; then
  go test ./... -covermode=count -coverprofile=coverage.out
else
  go test ./...
fi