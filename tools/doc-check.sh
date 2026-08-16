#!/usr/bin/env bash
# Doc link check - cross-document .md links must resolve.
#
# Separate responsibility from doc-lint.sh (markdown style): this guards link
# integrity (rot) as docs move, rename, or retire. Relative .md links only;
# http(s) and absolute paths are skipped, #anchors ignored.
#
# Cross-repo links are skipped, not failed: a target that escapes this repo root
# (a sibling repo, e.g. ../other-repo/...) or resolves into a git-ignored subtree
# (another repo's content mounted under this one, e.g. governance repo -> knowledge/)
# is another repo's responsibility and cannot be validated from here.
#
# Copy from docs/templates/markdownlint/ into tools/. Portable, local == CI.
set -uo pipefail
cd "$(cd "$(dirname "$0")/.." && pwd)" || exit 2   # repo root
root="$(pwd)"
in_git=0; git rev-parse --git-dir >/dev/null 2>&1 && in_git=1

# 대상 .md 목록 - git repo면 추적 파일만(=.gitignore·빌드/의존성 캐시 자동 제외),
# 아니면 find fallback. doc-lint의 gitignore:true와 스코프를 맞춘다.
list_md() {
  if git rev-parse --git-dir >/dev/null 2>&1; then
    git ls-files '*.md' ':!:**/legacy/**' ':!:CLAUDE.md' ':!:GEMINI.md'
  else
    find . -name '*.md' -not -type l \
      -not -path './node_modules/*' -not -path './.git/*' -not -path './legacy/*'
  fi
}

broken=0
while IFS= read -r f; do
  [ -z "$f" ] && continue
  dir=$(dirname "$f")
  while IFS= read -r link; do
    [ -z "$link" ] && continue
    case "$link" in http*|/*) continue ;; esac
    target=$(realpath -m "$dir/$link")
    # 크로스-repo 예외: 이 repo 밖(sibling)이거나 gitignore된 대상(다른 repo 콘텐츠)은 건너뛴다.
    case "$target" in "$root"/*) ;; *) continue ;; esac
    if [ "$in_git" = 1 ] && git check-ignore -q "$target" 2>/dev/null; then continue; fi
    [ -f "$target" ] || { printf '✗ %s -> %s\n' "$f" "$link"; broken=1; }
  done < <(grep -oP '(?<!\])\]\([^)#]+\.md' "$f" 2>/dev/null | sed -E 's/^\]\(//')
done < <(list_md)

if [ "$broken" -ne 0 ]; then
  echo "doc-check FAILED - broken .md links."
  exit 1
fi
echo "doc-check OK - all .md links resolve."
