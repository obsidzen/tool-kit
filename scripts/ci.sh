#!/usr/bin/env bash
set -euo pipefail

cd "$(dirname "${BASH_SOURCE[0]}")/.."

export PATH="$HOME/.local/bin:$PATH"
mise trust -a >/dev/null 2>&1 || true
mise install
eval "$(mise env -s bash)"

# 문서 게이트 — 래퍼가 vendored 돼 있으면 DoD의 일부 (docs/code-style.md §4).
if [ -f tools/doc-lint.sh ];  then bash tools/doc-lint.sh;  fi   # 마크다운 스타일(markdownlint)
if [ -f tools/doc-check.sh ]; then bash tools/doc-check.sh; fi   # 링크 무결성

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
