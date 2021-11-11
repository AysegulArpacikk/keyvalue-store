#!/bin/sh

cd "$(dirname "$0")"

echo ""
echo "GO Linter has been started"

golangci-lint run -c .golangci.yml -v --timeout 60s

echo ""
echo "GO Linter has been finished"
echo ""
