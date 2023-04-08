#!/bin/bash

go generate ./...
go test $(go list ./... | grep -v tools)