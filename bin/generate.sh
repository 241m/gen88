#!/usr/bin/env bash

cd "$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")" || exit 1

mkdir -p bin output

go build -o bin/gen88 main.go

for c in content/?*/text.txt; do
	name=$(basename "$(dirname "$c")")
	odir=output/$name
	mkdir -p "$odir"

	for s in symbols/?*/opt.svg; do
		n=$(basename "$(dirname "$s")")
		f=${odir}/${n}.svg

		echo "$f"

		bin/gen88 -f "$c" -s "$s" > "$f"
	done
done