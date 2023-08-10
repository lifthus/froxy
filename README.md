# Froxy

Simple proxy server which connects between local frontend and remote web server (especially development server) .

## Usage

- Starting proxy server on given port, which forwards all HTTP requests to given target origin.

```
froxy -p 8888 -t https://www.google.com
froxy --port 8888 --target https://www.google.com
```

- Generating goproxyconfig.json file

```
froxy
```
