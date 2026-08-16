#!/usr/bin/env bash
# Repo bootstrap (workspace convention). Run once per fresh clone.
# Activates the tracked git hooks and installs pinned tools. Idempotent.
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."

# 1) activate tracked hooks - core.hooksPath is local git config, so it is NOT
#    set automatically on clone.
if [ -d .githooks ]; then
  git config core.hooksPath .githooks
  echo "hooks: core.hooksPath -> .githooks"
fi

# 2) pinned toolchain
if command -v mise >/dev/null 2>&1; then
  mise install && echo "tools: mise install done"
else
  echo "tools: mise not found - skip (install mise for the pinned toolchain)"
fi

echo "setup OK"
