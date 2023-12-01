<p align="center">
<img src="https://github.com/lifthus/froxy/assets/108582413/de1db34c-6616-4745-9aa5-8000d34eb04d" width=300px />
</p>


# froxy

## Main features

HTTP, HTTPS 등의 프로토콜에 기반한 다양한 프록시 기능들을 손쉽게 사용할 수 있도록 해주는 개인용/테스트용 프록시 서버 애플리케이션.
다음 프록시 기능들을 froxyfile 설정을 통해 간단히 사용할 수 있다.

- 포워드 프록시
- 리버스 프록시
- 라운드-로빈 로드 밸런서

## Web dashboard

프록시 상태를 모니터링하고 관리할 수 있는 웹 대시보드가 내장되어 있다. 
웹 대시보드는 TLS 위에서 작동하도록 강제되며, 키 페어를 제공하지 않으면 self-signed certificate를 생성해 자동으로 TLS 기능을 상시 제공한다.
현재 다음과 같은 기능을 제공한다.

### Forward proxy
- 포워드 프록시 On/Off 스위치
- 포워드 프록시 화이트리스트(프록시를 사용할 수 있는 IP 주소 목록) 관리

### Reverse proxy
- 리버스 프록시 On/off 스위치
- 리버스 프록시 타겟 현황

## Security

키 페어를 제공하지 않아도 웹 대시보드는 자동으로 항상 HTTPS 상에서 동작한다.
리버스 프록시도 키 페어를 제공하지 않고 HTTPS 모드를 사용하면 self-signed certificate를 자동 생성한다.
이를 활용해 로컬 개발 환경에서 손쉽게 HTTPS를 적용해볼 수 있다.

## Installation

```
go install github.com/lifthus/froxy/cmd/froxy@latest
```

## Usage

커맨드를 실행하는 현재 디렉토리에 "froxyfile"이라는 이름의 설정 파일을 생성하고 다음 예시를 참고해서 작성한다.
이후 커맨드를 실행하면(그냥 쉘에 froxy 입력) froxyfile 설정에 따라 프록시 서버들이 설정되고 실행된다.

```yml
dashboard:
  port: 8542 # 포트를 따로 제공하지 않으면 기본 8542 포트에 대시보드가 실행된다.
  host: 123.123.123.123 # froxy 대시보드가 실행되는 서버의 IP 주소 혹은 도메인 네임.
  # 웹 대시보드는 TLS가 강제되며, 따라서 아래 "tls" 속성에 certificate와 key 파일이 제공되지 않으면,
  # froxy는 스스로 self-signed key 페어를 생성해 HTTPS 상에서 대시보드를 제공한다.
  tls:
    cert: ./cert.pem
    key: ./key.pem
# "forward"는 포워드 프록시 설정 리스트다.
forward:
  - name: forward-froxy # 해당 프록시를 식별하기 위한 고유한 프록시 서버 이름
    port: 8543 # 해당 프록시 서버가 실행될 포트 번호.
# "reverse"는 리버스 프록시 설정 리스트다.
reverse:
  - name: example-reverse
    port: 8544
    insecure: false # 명시적으로 true로 설정되면 리버스 프록시 프론트단은 아래 "tls" 속성이 존재하더라도 무시하고 HTTP 서버로 작동한다. 기본값은 false다.
    tls: # "insecure" 속성이 false임에도 불구하고 이 속성에서 키 페어 파일 경로가 주어지지 않으면, 대시보드의 경우 처럼 스스로 키 페어를 생성해 리버스 프록시 프론트단이 HTTPS 서버로 동작하도록 한다.
      cert: ./cert2.pem
      key: ./key2.pem
    proxy: # 리버스 프록시 경로 설정 부분.
      "abc.com": # 클라이언트 요청의 호스트 부분(타겟 호스트가 "abc.com"인 요청은 아래 설정에 따라 포워딩됨).
        "/": # "abc.com"으로의 모든 요청은 아래 URL로 포워딩된다.
          - http://127.0.0.1:8545
        "/api": # "/api"로 시작하는 모든 요청은 아래 URL로 포워딩된다.
          - http://127.0.0.1:8546
        "/api/v2": # "/api/v2"는 "/api"와 겹치지만 더 긴 base path를 가지는데, 이런 경우 더 긴 경로에 먼저 매칭된다.
          # 한 base path에 여러 타겟 URL이 제공되면 해당 타겟들에 대해 단순한 라운드-로빈 로드밸런서로 작동한다.
          - http://127.0.0.1:8547
          - http://127.0.0.1:8548

# Note: Base path "/api"에 대해 "http://cde.com/fgh"로 포워딩되는 경우, "/api/asdf"로 요청하면 "http://cde.com/fgh/asdf"로 포워딩된다.
```

사용자 로컬 환경에서는 물론 EC2 등 클라우드에 가상 컴퓨팅 환경을 가지고 있다면 설정 몇 줄과 커맨드 하나로 손쉽게 HTTP(S) 포워드 프록시를 구축할 수 있다.

로컬 개발 환경에서 간단히 리버스 프록시를 설정해 개발 서버와의 쿠키 연동 문제 같은 다양한 문제들을 손쉽게 해결할 수 있다(프론트엔드 개발 서버에서 보통 프록시 기능을 제공하긴 하지만...).

예컨대 froxy는, froxy 자체를 이용해 path를 기반으로 리버스 프록시를 설정해 한쪽은 리액트 개발 서버로, 한쪽은 froxy 자신의 대시보드 API로 포워딩함으로써 하나의 서버처럼 동작하도록 하여 개발을 진행했다.

## Getting started

먼저 Go 언어 환경을 갖추고 위 Installation에 기술된 커맨드를 통해 froxy를 설치한 후, 위 예시를 참고해 필요에 따라 froxyfile을 작성한다.

### Execution
#### Foreground
1. 셸을 계속 켜놓고 포그라운드에서만 실행하려면 간단히 "froxy" 커맨드만 입력하고 계속 진행하면 된다.
#### Background
1. 셸을 종료하고도 백 그라운드에서 계속 실행하려면, "froxy" 커맨드를 실행한 후, 대시보드에서 사용할 루트 계정 정보를 입력하고, Ctrl + Z로 빠져나온다.
2. "bg" 커맨드를 통해 간단히 지금 일시정지된 froxy를 백그라운드에서 계속 실행하도록 한다.
3. "disown" 커맨드를 통해 간단히 셸이 정지되고 나서도 froxy가 계속 백그라운드에서 실행되도록 한다.

### Termination
Foreground의 경우 Ctrl+C로 빠져나온다.

Background의 경우 "pgrep froxy" 커맨드를 통해 froxy의 PID를 찾은 후 "kill [PID]" 커맨드를 통해 프로세스를 종료한다.

### Dashboard
1. https://[대시보드 호스트]:[포트]로 접속해, 다음 화면에서 froxy를 실행하며 설정한 루트 계정으로 로그인한다.

<p align="center">
<img width="274" alt="image" src="https://github.com/lifthus/froxy/assets/108582413/e5e99069-432e-45c6-8ab7-967de836522d">
</p>

2. 로그인하면 다음과 같이 설정된 각 서버들의 상태를 간단히 볼 수 있다. 왼쪽 파란 버튼을 통해 각 서버를 끄거나 켤 수 있다.
<p align="center">
<img width="257" alt="image" src="https://github.com/lifthus/froxy/assets/108582413/c30fa6f3-ace6-4d37-b57d-b34b62ef9ec1">
</p>

3. 상단 상태바 최우측에 있는 주황색 버튼을 통해 로그아웃 할 수 있으며, 바로 그 왼쪽에서 대시보드에 접속한 자신의 IP 주소를 확인할 수 있다.
<p align="center">
<img width="169" alt="image" src="https://github.com/lifthus/froxy/assets/108582413/882f7cd0-8fb7-49bf-bff3-c738acac23f6">
</p>

#### Forward proxy
메인 화면에서 포워드 프록시 하나를 클릭해서 들어가면 다음과 같이 포워드 프록시를 사용할 수 있는 IP 주소 목록을 관리할 수 있다.

상단 상태바에서 자신의 IP 주소를 확인 후 화이트리스트에 추가하고 해당 디바이스 브라우저나 네트워크 설정에서 포워드 프록시 호스트와 포트를 설정하면 모든 HTTP(S) 트래픽을 해당 서버로 우회시킬 수 있다.

프록시 이름과 포트 번호 아래의 큰 버튼을 통해 프록시를 끄고 켤 수 있다.
<p align="center">
<img width="268" alt="image" src="https://github.com/lifthus/froxy/assets/108582413/e5d48c88-8352-41f1-8d84-74687849da7a">
</p>

#### Reverse proxy
메인 화면에서 리버스 프록시 하나를 클릭해서 들어가면 다음과 같이 리버스 프록시의 포워딩 테이블을 확인할 수 있다.

프록시 이름과 포트 번호 아래의 큰 버튼을 통해 프록시를 끄고 켤 수 있다.
<p align="center">
  <img width="496" alt="image" src="https://github.com/lifthus/froxy/assets/108582413/878e59a5-5bdb-40ed-9c8c-c1e22aee0756">
</p>





