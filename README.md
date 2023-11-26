# froxy

## Main features

HTTP, HTTPS 등의 프로토콜에 기반한 다양한 프록시 기능들을 손쉽게 사용할 수 있도록 해주는 개인용 프록시 서버 애플리케이션.
다음 기능들을 froxyfile 설정을 통해 간단히 사용할 수 있다.

- 포워드 프록시
- 리버스 프록시
- 로드 밸런서

## Web dashboard

프록시 상태를 모니터링하고 관리할 수 있는 웹 대시보드가 내장되어 있다. 
웹 대시보드는 TLS 위에서 작동하도록 강제되며, 키 페어를 제공하지 않으면 self-signed certificate를 생성해 자동으로 TLS 기능을 상시 제공한다.
현재 다음과 같은 기능을 제공한다.

- 포워드 프록시 화이트리스트(프록시를 사용할 수 있는 IP 주소 목록) 설정

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
이후 커맨드를 실행하면 설정에 따라 프록시 서버들이 설정되고 실행된다.

```yml
dashboard:
  port: 8542 # 8542 is default port for froxy dashboard.
  host: 123.123.123.123 # dashboard host as an IP addr or domain name
  # HTTPS is mandatory for web dashboard, so if "tls" not set,
  # froxy automatically generates self-signed key pair for HTTPS.
  tls:
    cert: ./cert.pem
    key: ./key.pem
# forward is config list of forward proxies.
forward:
  - name: forward-froxy # name to identify the proxy
    port: 8543 # port to listen on
# reverse is config list of reverse proxies.
reverse:
  - name: example-reverse
    port: 8544
    insecure: false # if explcitly set to true, HTTP is used instead of HTTPS(ignoring "tls" field). Default is false.
    tls: # if not set when insecure flag is false, self-signed certificates for given hosts are automatically generated.
      cert: ./cert2.pem
      key: ./key2.pem
    proxy:
      "abc.com": # target host of request message to match. "*" can be used to match all kinds of hosts.
        "/": # all requests to abc.com will be forwarded to the URL specified in this list.
          - http://127.0.0.1:8545
        "/api": # all request with target path starting with "/api" will be forwarded to the url specified in this list.
          - http://127.0.0.1:8546
        "/api/v2": # note that longer path(/api/v2) will be matched first, prior to shorter path(/api).
          # if multiple URLs provided, the reverse proxy automatically works as a round-robin load balancer.
          - http://127.0.0.1:8547
          - http://127.0.0.1:8548
```

사용자 로컬 환경에서는 물론 EC2 등 클라우드에 가상 컴퓨팅 환경을 가지고 있다면 설정 몇 줄과 커맨드 하나로 손쉽게 HTTP(S) 포워드 프록시를 구축할 수 있다.

로컬 개발 환경에서 간단히 리버스 프록시를 설정해 개발 서버와의 쿠키 연동 문제 같은 다양한 문제들을 손쉽게 해결할 수 있다(프론트엔드 개발 서버에서 보통 프록시 기능을 제공하긴 하지만...).

예컨대 froxy는, froxy 자체를 이용해 path를 기반으로 리버스 프록시를 설정해 한쪽은 리액트 개발 서버로, 한쪽은 froxy 자신의 대시보드 API로 포워딩함으로써 하나의 서버처럼 동작하도록 하여 개발을 진행했다.
