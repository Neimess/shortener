#!/bin/bash
set -a
source "${WORKDIR}/.env.local"
set +a
go run "${WORKDIR}/cmd/shortener/main.go"
