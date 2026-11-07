# commit-msg — 커밋 메시지 검사

## 만든 이유

공통 커밋 규약을 적용하되 제목 검사 정규식을 직접 유지하지 않기 위해 도입했다.

## 역할과 범위

committed가 Conventional Commits를 검사하고 prek가 AI 작성 푸터 패턴을 검사한다. 자동 생성 Merge·Revert·fixup 제목에 예외를 두지 않는다. 실제 merge는 conventional 제목으로 수행할 수 있다.

## 구현과 연결

[진입점](../../.githooks/commit-msg), [검사 설정](../../.pre-commit-config.yaml), [committed 설정](../../committed.toml)을 사용한다.

## 검증과 한계

정상·비정상 제목과 AI 푸터를 임시 메시지 파일로 확인했다. AI 푸터 검사는 메시지 전체의 지정 패턴을 확인하며 AI 작성 여부를 판별하는 기능은 아니다.
