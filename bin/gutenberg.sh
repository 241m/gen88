#!/usr/bin/env bash

if [ $# -lt 2 ] || [ "$1" = -h ] || [ "$1" = --help ]; then
	echo "Download and clean up plain text file from Project Gutenberg."
	echo "Url must be a link to the work's plain text UTF-8 format. "
	echo
	echo "Usage: $(basename "$0") <name> <url>"
	exit 1
fi

cd "$(realpath "$(dirname "${BASH_SOURCE[0]}")/..")" || exit 1

name=${1:?}
url=${2:?}
odir=content/$name

mkdir -p "$odir"
curl  -o "$odir/src.txt" -L "$url" | sed -n '/\*\*\* START/,/\*\*\* END/{//!p;}'