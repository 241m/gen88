#!/usr/bin/env bash

cd "$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")" || exit 1

data=$(date +%s)

for s in symbols/?*/opt.svg; do
	n=$(basename "$(dirname "$s")")
	f="output/${n}-${data}.svg"
	echo "$f"
	echo "$data" | go run main.go -f - -s "$s" | xmllint --format - > "$f"
done