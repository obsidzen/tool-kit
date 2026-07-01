#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

export PATH="$HOME/.local/bin:$PATH"
mise trust -a >/dev/null 2>&1 || true
mise install
eval "$(mise env -s bash)"

export GOCACHE="${GOCACHE:-$(pwd)/.cache/go-build}"
export GOFLAGS="${GOFLAGS:--buildvcs=false}"
mkdir -p "$GOCACHE"

go version

modules=(
  cli-kit
  run-kit
  tui-kit
)

for m in "${modules[@]}"; do
  echo "== $m =="
  ( cd "$m" && go build ./... && go vet ./... && go test ./... )
done

echo "tool-kit CI OK"
