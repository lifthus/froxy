# froxy

HTTP, HTTPS 상에서 다양한 기능을 제공하는 프록시 서버 애플리케이션.
다음 기능들을 커맨드와 옵션 플래그 한 두 개로 간단히 사용할 수 있다.

- 포워드 프록시
- 리버스 프록시
- 로드 밸런서

## Installation

```
  go install github.com/lifthus/froxy/cmd/froxy@latest
```

## Usage

커맨드를 실행하는 현재 디렉토리에 "froxyfile"이라는 이름의 설정 파일을 다음 예시를 참고해서 작성한다.
이후 커맨드를 실행하면 설정에 따라 프록시 서버들이 설정되고 실행된다.

```yml
# forward is config list of forward proxies.
forward:
  - name: forward-froxy # name to identify the proxy
    port: 8543 # port to listen on
    # allowed is a list of IP addrs allowed to use the forward proxy, without authentication.
    # You can allow some IP addrs to use the proxy without enabling web dashboard.
    # note that localhost is always allowed.
    allowed:
      - "*" # all IP addrs are allowed
      - 123.123.123.123
# reverse is config list of reverse proxies.
reverse: # reverse proxy config
  - name: example-reverse
    port: 8544
    insecure: false # if explcitly set to true, HTTP is used instead of HTTPS(ignoring "tls" field of each "proxy"). Default is false.
    proxy:
      - host: abc.com # target host of request message to match (optional. in default, host is simply ignored).
        tls: # if not set, self-signed certificate for given host is automatically generated.
          cert: ./cert2.pem
          key: ./key2.pem
        target: # proxy config
          - path: / # all requests to abc.com will be forwarded to the url specified in "to" property below.
            to:
              - http://127.0.0.1:8545
          - path: /api # all request with target url starting with "/api" will be forwarded to the url specified in "to" property below.
            to:
              - http://127.0.0.1:8546
          - path: /api/v2 # note that longer path(/api/v2) will be matched first, prior to shorter path(/api).
            to: # if you provide multiple urls, the reverse proxy automatically works as a round-robin load balancer.
              - http://127.0.0.1:8547
              - http://127.0.0.1:8548
```
