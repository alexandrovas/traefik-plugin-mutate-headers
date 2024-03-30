# Mutate Headers

Mutate Headers is a middleware plugin for [Traefik](https://traefik.io) which mutates the request headers based on the provided configuration.

## Configuration

### Fields

- `header` (string, required): The name of the header to be mutated.
- `newName` (string, optional): The new name of the header. If not provided, the header will be mutated in place.
- `regex` (string, optional): The regular expression to match the header value. If not set the header value will be preserved.
- `replacement` (string, optional): The replacement string for the header value. Must be set if `regex` is set.
- `deleteSource` (bool, optional): If set to true, the source header will be deleted after the mutation. Useful when renaming headers.

### Static

```yaml
experimental:
  plugins:
    traefik-plugin-mutate-headers:
      modulename: "github.com/trolleksii/traefik-plugin-mutate-headers"
      version: "v0.1.2"
```

### Dynamic

To configure the Mutate Headers plugin you should create a [middleware](https://docs.traefik.io/middlewares/overview/) in your dynamic configuration as explained [here](https://docs.traefik.io/middlewares/overview/). 
The following example creates and uses the mutateHeaders middleware plugin to set the `X-Host` header to the particular subdomain of the `Host` header.

```yaml
http:
  routes:
    my-router:
      rule: "Host(`localhost`)"
      service: "my-service"
      middlewares : 
        - "mutateHeaders"
  services:
    my-service:
      loadBalancer:
        servers:
          - url: "http://127.0.0.1"
  middlewares:
    rewriteHeaders:
      plugin:
        mutateHeaders:
          mutations:
            - header: "Host"
              newName: "X-Host"
              regex: "^(.+)\.testing\.com$"
              replacement: "$1"
```
