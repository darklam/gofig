#!/usr/bin/env bash

docker run --rm -v "$PWD":/src -w /src vektra/mockery --all