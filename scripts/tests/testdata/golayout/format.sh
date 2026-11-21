#!/usr/bin/env bash
set -euo pipefail
[[ "$1" == setup ]]
[[ "$2/tools" == "$3" ]]
[[ -z "$GOFLAGS" ]]
[[ "$GOWORK" != off ]]
[[ "$("$GO" -C "$3" list -m)" == github.com/longxi-go/gosr/tools ]]
[[ "$("$GO" -C "$3" list -m -f '{{.Dir}}' github.com/longxi-go/gosr)" == "$2" ]]
printf '%s\n' "$4" > "$GOLAYOUT_TEST_LANGUAGE"
