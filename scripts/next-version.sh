#!/usr/bin/env bash
# Print the next v2 patch tag (v2.0.0 when no prior v2 tag exists).
set -euo pipefail

latest="$(git tag -l 'v2.*' --sort=-v:refname | head -1 || true)"

if [[ -z "$latest" ]]; then
  printf 'v2.0.0'
  exit 0
fi

ver="${latest#v}"
IFS=. read -r major minor patch _ <<< "$ver"
patch=$((patch + 1))
printf 'v%s.%s.%s' "$major" "$minor" "$patch"