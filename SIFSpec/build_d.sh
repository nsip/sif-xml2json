#!/bin/bash
set -e

CGO_ENABLED=0 go run ./1_txt2toml/main.go
CGO_ENABLED=0 go run ./2_toml2json/config.go ./2_toml2json/main.go
