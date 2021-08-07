# Traefik CORS Plugin

Automatically handle preflight requests with Traefik.

# Usage

This plugin creates a Traefik [Middleware](https://doc.traefik.io/traefik/middlewares/overview/) to handle CORS headers. To use any Traefik plugin, you must first enable it in the [static configuration](https://doc.traefik.io/traefik/getting-started/configuration-overview/#the-static-configuration):

```yaml
experimental:
  plugins:
    cors:
      modulename: git.quintin.dev/networking/traefik-cors
      version: 0.1.0
```

Then it can be used in the [dynamic configuration](https://doc.traefik.io/traefik/getting-started/configuration-overview/#the-dynamic-configuration). Here is an example with the [Kubernetes CRD Provider](https://doc.traefik.io/traefik/reference/dynamic-configuration/kubernetes-crd/):

```yaml
apiVersion: traefik.containo.us/v1alpha1
kind: Middleware
metadata:
  name: cors
spec:
  cors:
    AllowCredentials: false
    AllowHeaders: []
    AllowMethods:
    - GET
    - HEAD
    - POST
    AllowOrigins:
    - "*"
    ExposeHeaders: []
    MaxAge: 5
```

The values shown in the above `Middleware` are all the default values. Applying an empty middleware of type `cors` to your route will result in these values.

### `AllowCredentials`

Configures the [Access-Control-Allow-Credentials](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Credentials) header.

Weather or not clients may use credentials mode `"include"`. When clients do use credentials mode `"include"`, all wildcards in other headers are treated as string literals.

> Note: If you need credentials from a client (which includes cookies, TLS client certificates, and authentication entries), this needs to be enabled.

### `AllowHeaders`

Configures the [Access-Control-Allow-Headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Headers) header.

The list of headers to allow from clients. If `"*"` is present, the wildcard value will _always_ be returned. If clients use credentials mode `"include"`, the wildcard value is treated as the literal string `"*"`. The `Authorization` header is not included in the wildcard.

> Note: If you need credentials from a client or the `Authorization` header, you cannot use wildcard (`"*"`).

### `AllowMethods`

Configures the [Access-Control-Allow-Methods](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Methods) header.

The list of methods to allow from clients. If `"*"` is present, the wildcard value will _always_ be returned. If clients use credentials mode `"include"`, the wildcard value is treated as the literal string `"*"`.

> Note: If you need credentials from a client, you cannot use wildcard (`"*"`).

### `AllowOrigins`

Configures the [Access-Control-Allow-Origin](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Allow-Origin) header.

The list of origins to allow from clients. If `"*"` is present, the wildcard value will _always_ be returned. If `"*"` is present, no other values should be provided. While this will still work, it is less efficient for both the client and server. If clients use credentials mode `"include"`, the wildcard value will cause clients to fail.

The origin `"null"` cannot be configured via this plugin and [should not be used](https://w3c.github.io/webappsec-cors-for-developers/#avoid-returning-access-control-allow-origin-null).

> Note: If you need credentials from a client, you cannot use wildcard (`"*"`).

### `ExposeHeaders`

Configures the [Access-Control-Expose-Headers](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Expose-Headers) header.

The list of headers to expose from the server. If `"*"` is present, the wildcard value will _always_ be returned. If clients use credentials mode `"include"`, the wildcard value is treated as the literal string `"*"`. The `Authorization` header is not included in the wildcard.

> Note: If you need credentials from a client or the `Authorization` header, you cannot use wildcard (`"*"`).

### `MaxAge`

Configures the [Access-Control-Max-Age](https://developer.mozilla.org/en-US/docs/Web/HTTP/Headers/Access-Control-Max-Age) header.

The duration of time in seconds that a CORS preflight request should be cached for. This cannot exceed the browser defined maximum. Many browsers will cache for 5 seconds if this header is not included. `-1` can be used to disable caching of preflight requests.

# FAQ's

### Doesn't Traefik already handle CORS?

Traefik has an amazing [headers](https://doc.traefik.io/traefik/middlewares/headers/) middleware that allows routes to append headers to every response. This middleware even comes with support for CORS headers already! The only issue with this middleware is that it does not end processing on CORS preflight requests after writing the headers. This means the backend application has to be prepared to accept `OPTIONS` requests, or the backend service needs to be routed to a no-op on preflight requests. Neither of these situations is ideal. This plugin solves this issues by stopping the preflight request and returning immediately, without contacting the backend application.

# Resources

#### [Fetch Specification - CORS Protocol](https://fetch.spec.whatwg.org/#http-cors-protocol)

#### [MDN - Cross-Origin Resource Sharing](https://developer.mozilla.org/en-US/docs/Web/HTTP/CORS)

#### [W3C - CORS for Developers](https://w3c.github.io/webappsec-cors-for-developers/)
