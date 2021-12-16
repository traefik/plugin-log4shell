# Log4Shell Mitigation

[![Build Status](https://github.com/traefik/plugin-log4shell/workflows/Main/badge.svg?branch=master)](https://github.com/traefik/plugin-log4shell/actions)

Log4Shell is a middleware plugin for [Traefik](https://github.com/traefik/traefik) which blocks JNDI attacks based on HTTP header values.

Related to the Log4J CVE: https://cve.mitre.org/cgi-bin/cvename.cgi?name=CVE-2021-44228

## Configuration

### Static

```bash
--pilot.token=xxx
--experimental.plugins.log4shell.modulename=github.com/traefik/plugin-log4shell
--experimental.plugins.log4shell.version=v0.1.2
```

```yaml
pilot:
  token: xxx

experimental:
  plugins:
    log4shell:
      modulename: github.com/traefik/plugin-log4shell
      version: v0.1.2
```

```toml
[pilot]
    token = "xxx"

[experimental.plugins.log4shell]
    modulename = "github.com/traefik/plugin-log4shell"
    version = "v0.1.2"
```

### Dynamic

To configure the `Log4Shell` plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/).

#### File 

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

#### Kubernetes

```yaml
---
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: log4shell-foo
spec:
  plugin:
    log4shell:
      errorCode: 200

---
apiVersion: traefik.containo.us/v1alpha1
kind: IngressRoute
metadata:
  name: whoami
spec:
  entryPoints:
    - web
  routes:
    - kind: Rule
      match: Host(`whoami.example.com`)
      middlewares:
        - name: log4shell-foo
      services:
        - kind: Service
          name: whoami-svc
          port: 80
```

```yaml
---
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: log4shell-foo
spec:
  plugin:
    log4shell:
      errorCode: secretName

---
apiVersion: networking.k8s.io/v1
kind: Ingress
metadata:
  name: myingress
  annotations:
    traefik.ingress.kubernetes.io/router.middlewares: default-log4shell-foo@kubernetescrd

spec:
  rules:
    - host: example.com
      http:
        paths:
          - path: /
            pathType: Prefix
            backend:
              service:
                name:  whoami
                port:
                  number: 80
```

#### Docker

```yaml
version: '3.7'

services:
  whoami:
    image: traefik/whoami:v1.7.1
    labels:
      traefik.enable: 'true'

      traefik.http.routers.app.rule: Host(`whoami.localhost`)
      traefik.http.routers.app.entrypoints: websecure
      traefik.http.routers.app.middlewares: log4shell-foo
      
      traefik.http.middlewares.log4shell-foo.plugin.log4shell.errorcode: 200
```
