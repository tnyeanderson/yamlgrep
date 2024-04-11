#!/bin/bash

filter() {
	if grep "$@" >/dev/null 2>&1; then
		echo "$doc"
	fi
}

doc=''
while IFS= read -r line; do
	if [[ "${line:0:3}" == '---' ]] && [[ -n "$doc" ]]; then
		filter "$@" <<<"$doc"
		doc=''
	fi
	if [[ -n "$doc" ]]; then
		doc+=$'\n'
	fi
	doc+="$line"
done
filter "$@" <<<"$doc"
