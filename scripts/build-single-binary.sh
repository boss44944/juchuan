#!/bin/sh
set -eu

root=$(CDPATH='' cd -- "$(dirname -- "$0")/.." && pwd)
output=${1:-"$root/dist/juchuan"}
case "$output" in
  /*) ;;
  *) output="$PWD/$output" ;;
esac

stage=$(mktemp -d "${TMPDIR:-/tmp}/juchuan-build.XXXXXX")
trap 'rm -rf -- "$stage"' EXIT
trap 'exit 1' HUP INT TERM

mkdir -p "$stage/frontend" "$stage/backend" "$(dirname -- "$output")"
cp "$root/frontend/package.json" "$root/frontend/package-lock.json" \
  "$root/frontend/index.html" "$root/frontend/vite.config.ts" "$stage/frontend/"
cp -R "$root/frontend/src" "$root/frontend/public" "$stage/frontend/"
cp "$root"/backend/*.go "$root/backend/go.mod" "$root/backend/go.sum" \
  "$root/backend/app-logo.png" "$stage/backend/"

(
  cd "$stage/frontend"
  npm ci --ignore-scripts --no-audit --no-fund
  npm run build
)
mv "$stage/frontend/dist" "$stage/backend/static"
(
  cd "$stage/backend"
  go build -trimpath -ldflags '-X main.frontendEmbeddedOnly=true' -o "$output" .
)
printf 'Built standalone Juchuan: %s\n' "$output"
