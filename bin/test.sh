#!/usr/bin/env bash

cd "$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")" || exit 1

f=output/test-$(date +%s).svg

go run main.go "$@" | xmllint --format - > "$f"

echo "$f"