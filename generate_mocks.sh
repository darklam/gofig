#!/bin/bash

moq -out ./mocks/provider.go -pkg mocks ./interfaces Provider
moq -out ./mocks/valuer.go -pkg mocks ./interfaces Valuer