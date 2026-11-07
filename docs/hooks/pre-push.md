# pre-push — push 전 전체 검증

## 만든 이유

원격 전송 전에 저장소 전체 검증을 수행하면서 커밋 검사를 가볍게 유지하기 위해 도입했다.

## 역할과 범위

prek가 기존 scripts/ci.sh를 호출한다. 문서 검사와 cli-kit·run-kit·tui-kit의 build·vet·test가 실패하면 push를 중단한다. 작업 완료 조건의 전체 CI는 유지한다.

## 구현과 연결

[진입점](../../.githooks/pre-push), [검사 설정](../../.pre-commit-config.yaml), [전체 CI](../../scripts/ci.sh)를 사용한다. 서버 CI도 같은 스크립트를 호출한다.

## 검증과 한계

실패 CI의 차단을 임시 저장소에서 확인했다. 실제 원격 push는 실행하지 않는다. 이 저장소의 속성에는 LFS 필터가 없어 LFS 의존성과 원격 이름 제한을 추가하지 않았다. 로컬 hook은 우회할 수 있으므로 서버 검증을 유지한다.
