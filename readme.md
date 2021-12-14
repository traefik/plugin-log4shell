# Log4Shell Mitigation

[![Build Status](https://github.com/traefik/plugin-log4shell/workflows/Main/badge.svg?branch=master)](https://github.com/traefik/plugin-log4shell/actions)

Log4Shell is a middleware plugin for [Traefik](https://github.com/traefik/traefik) which blocks JNDI attacks based on HTTP header values.

Related to the Log4J CVE: https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-44228

## Configuration

## Static

```yaml
pilot:
  token: xxx

experimental:
  plugins:
    log4shell:
      modulename: github.com/traefik/plugin-log4shell
      version: v0.1.0
```

```toml
[pilot]
    token = "xxx"

[experimental.plugins.log4shell]
    modulename = "github.com/traefik/plugin-log4shell"
    version = "v0.1.0"
```

## Dynamic

To configure the `Log4Shell` plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/).

```yaml
http:
  middlewares:
    log4shell-foo:
      plugin:
        log4shell:
          errorCode: 200

  routers:
    my-router:
      rule: Host(`localhost`)
      middlewares:
        - log4shell-foo
      service: my-service

  services:
    my-service:
      loadBalancer:
        servers:
          - url: 'http://127.0.0.1'
```

```toml
[http.middlewares]
  [http.middlewares.log4shell-foo.plugin.log4shell]
    errorCode = 200

[http.routers]
  [http.routers.my-router]
    rule = "Host(`localhost`)"
    middlewares = ["log4shell-foo"]
    service = "my-service"

[http.services]
  [http.services.my-service]
    [http.services.my-service.loadBalancer]
      [[http.services.my-service.loadBalancer.servers]]
        url = "http://127.0.0.1"
```
