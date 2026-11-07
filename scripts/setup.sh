#!/usr/bin/env bash
set -euo pipefail
cd "$(dirname "${BASH_SOURCE[0]}")/.."

command -v mise >/dev/null 2>&1 || { echo "setup: mise is required." >&2; exit 2; }
mise install
mise exec -- prek validate-config .pre-commit-config.yaml
git config core.hooksPath .githooks
printf '%s\n' "setup OK: fast commit checks and full pre-push CI are connected."
