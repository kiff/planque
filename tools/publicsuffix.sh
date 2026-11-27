#!/usr/bin/env bash
# Refreshes internal/browser/public_suffix_list.dat from the Public Suffix
# List (MPL 2.0): its rules, with comments removed and the license kept.
# Review the diff like any data change.
set -euo pipefail
out=$(cd "$(dirname "$0")/.." && pwd)/internal/browser/public_suffix_list.dat
tmp=$(mktemp)
trap 'rm -f "$tmp"' EXIT
curl -fsS --max-time 60 -o "$tmp" https://publicsuffix.org/list/public_suffix_list.dat
{
  printf '%s\n' "// This Source Code Form is subject to the terms of the Mozilla Public" \
    "// License, v. 2.0. If a copy of the MPL was not distributed with this" \
    "// file, You can obtain one at https://mozilla.org/MPL/2.0/." "//" \
    "// The Public Suffix List's rules (https://publicsuffix.org/list/public_suffix_list.dat)," \
    "// fetched $(date -u +%F), comments removed; refresh with tools/publicsuffix.sh." ""
  awk '/^\/\/ ===(BEGIN|END)/ { print; next } /^\/\// || /^[[:space:]]*$/ { next } { print $1 }' "$tmp"
} >"$out"
echo "wrote $out ($(grep -vc '^//' "$out") rules)"
