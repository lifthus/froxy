# Froxy

Simple proxy server which connects between localhost and remote web server (especially for frontend and development server).
It runs a proxy server on given port, forwarding all HTTP requests to given target origin.

## Usage

- Starting proxy server on given port, which forwards all HTTP requests to given target origin.

```
froxy -p 8888 -t https://www.google.com
froxy --port 8888 --target https://www.google.com
```

- Generating **froxyconfig.json** file
- Starting with **froxyconfig.json**

```
# generate froxyconfig.json template
froxy
# set the configuration
# and start froxy
froxy
```
