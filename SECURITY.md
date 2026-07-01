# 보안 정책

[한국어](SECURITY.md) · [English](SECURITY.en.md)

## 신고 방법

보안 취약점은 public issue, discussion, PR에 올리지 않는다.

가능하면 GitHub의 private vulnerability reporting을 사용한다. 사용할 수 없다면 maintainer에게 비공개 채널로 연락한다.

신고에는 가능한 범위에서 다음 정보를 포함한다.

- 영향받는 module: `cli-kit`, `run-kit`, `tui-kit` 중 하나
- 영향받는 version 또는 commit
- 재현 절차
- 예상 영향
- 관련 로그, 설정, 입력값
- 공개해도 되는 최소 proof of concept

## 범위

지원 범위:

- 현재 `main`
- 최신 정식 module release
- 최근 release candidate, 있는 경우

명시적으로 지원하지 않는 범위:

- fork에서만 존재하는 변경
- 지원 종료된 release line
- 소비 tool의 자체 코드나 운영 환경 문제
- 로컬 설정 실수, 노출된 개인 secret, 또는 사용자의 운영 환경 자체 문제

## 기대 응답

- 신고를 받으면 가능한 한 빨리 접수 여부를 알린다.
- 재현 가능성과 영향을 확인한 뒤 수정 방향을 정한다.
- 필요한 경우 보안 수정 release와 advisory를 함께 준비한다.

응답 시간은 보장하지 않지만, public issue로 먼저 공개하는 것은 피한다.

## Secret 처리

- token, key, `.env` 실값, 개인 endpoint는 신고 본문에 직접 붙이지 않는다.
- 이미 노출된 secret은 즉시 폐기하고 새 값으로 교체한다.
- secret 노출이 의심되면 code fix와 별도로 회전(rotation) 여부를 확인한다.

## 공개 원칙

취약점은 수정안, release, 사용자 완화책이 준비된 뒤 공개한다. 악용 가능성이 큰 세부 정보는 필요한 범위로 제한한다.
