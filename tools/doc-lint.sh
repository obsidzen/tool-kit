#!/usr/bin/env bash
# Doc lint - Markdown consistency via markdownlint-cli2 (DavidAnson).
#
# Standard markdownlint hygiene + house custom rules (tools/markdownlint/:
# no-emoji, no-circled), configured by .markdownlint-cli2.jsonc. Link integrity
# is a separate gate: tools/doc-check.sh.
#
# Copy this bundle from docs/templates/markdownlint/ into a repo:
#   markdownlint-cli2.jsonc -> .markdownlint-cli2.jsonc
#   rules/*.mjs             -> tools/markdownlint/
#   doc-lint.sh doc-check.sh -> tools/
# Portable, local == CI. Needs node + npx (markdownlint-cli2 is fetched/cached).
set -uo pipefail
cd "$(cd "$(dirname "$0")/.." && pwd)" || exit 2   # repo root

exec npx --yes markdownlint-cli2 "**/*.md"
