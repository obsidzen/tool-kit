# 기여 가이드

[한국어](CONTRIBUTING.md) · [English](CONTRIBUTING.en.md)

이 문서는 `tool-kit`에 기여할 때 따르는 사람 대상 규칙이다. AI 에이전트 작업 규칙은 `AGENTS.md`를 따른다.

## 개발 흐름

- GitHub-flow를 사용한다. `main`에 직접 push하지 않고 짧은 작업 브랜치에서 PR을 만든다.
- PR 제목과 commit message는 Conventional Commits 형식을 따른다: `<type>(<scope>): <subject>`.
- 기본 merge 방식은 squash merge다.

## 시작하기

```bash
./scripts/ci.sh
```

## Definition of Done

PR을 열기 전에 다음을 확인한다.

- [ ] `./scripts/ci.sh`가 통과한다.
- [ ] public API, import path, CLI/TUI 사용법이 바뀌면 해당 module의 `README.md`와 `README.en.md`를 함께 갱신했다.
- [ ] repo 운영 규칙이나 보안 신고 방식이 바뀌면 `CONTRIBUTING.md`/`CONTRIBUTING.en.md` 또는 `SECURITY.md`/`SECURITY.en.md`를 함께 갱신했다.
- [ ] 새 동작에는 필요한 테스트나 검증 절차가 있다.

## 코딩 규칙

- 각 kit은 작고 독립적인 Go module로 유지한다.
- 소비 tool이 Bubble Tea/Bubbles를 직접 의존해야 하는 경우를 먼저 `tui-kit` adapter나 reusable model 후보로 검토한다.
- 외부 프로세스 실행은 테스트 가능한 `run-kit.Runner` 경계 뒤에 둔다.
- 공개 API는 작고 명시적으로 유지한다.
- 주석은 코드가 말하지 못하는 이유, 불변식, 주의점을 설명할 때만 쓴다.

## 이슈와 PR

- 버그, 기능 요청, 작업, 질문은 `.github/ISSUE_TEMPLATE/`의 폼을 사용한다.
- PR 본문은 `.github/pull_request_template.md`를 따른다.
- 보안 취약점은 public issue나 PR에 쓰지 않는다. `SECURITY.md`를 따른다.

## 릴리스

- 각 module은 독립적으로 version pinning한다.
- release tag는 Go multi-module 규칙에 맞춰 `<module>/vX.Y.Z` 형식을 사용한다.
- tag는 `./scripts/ci.sh`를 통과한 commit에만 단다.
