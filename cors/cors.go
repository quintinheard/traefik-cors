// Package cors implements the CORS operations defined in the fetch spec.
// https://fetch.spec.whatwg.org/#http-cors-protocol
package cors

import (
	"net/http"
	"strconv"
	"strings"
)

const (
	// HeaderOrigin request header indicates where a fetch originates from.
	// See: RFC6454 § 7. The HTTP Origin Header Field
	// See: Fetch Standard § 3.1. `Origin` header.
	HeaderOrigin = "Origin"
	// HeaderVary indicates clients should not cache response headers.
	// See: RFC7231 § 7.1.4. Vary.
	// See: Fetch Standard § CORS protocol and HTTP caches.
	HeaderVary = "Vary"

	// HeaderAllowCredentials indicates whether the response can be shared when
	// request’s credentials mode is "include".
	// See: Fetch Standard § 3.2.3. HTTP responses.
	HeaderAllowCredentials = "Access-Control-Allow-Credentials"
	// HeaderAllowOrigin indicates whether the response can be shared, via returning
	// the literal value of the `Origin` request header (which can be `null`) or `*`
	// in a response.
	// See: Fetch Standard § 3.2.3. HTTP responses.
	HeaderAllowOrigin = "Access-Control-Allow-Origin"
	// HeaderAllowHeaders indicates which headers are supported by the response’s URL
	// for the purposes of the CORS protocol.
	// See: Fetch Standard § 3.2.3. HTTP responses.
	HeaderAllowHeaders = "Access-Control-Allow-Headers"
	// HeaderAllowMethods indicates which methods are supported by the response’s URL for
	// the purposes of the CORS protocol.
	// See: Fetch Standard § 3.2.3. HTTP responses.
	HeaderAllowMethods = "Access-Control-Allow-Methods"
	// HeaderExposeHeaders indicates which headers can be exposed as part of the response by listing their names.
	// See: Fetch Standard § 3.2.3. HTTP responses.
	HeaderExposeHeaders = "Access-Control-Expose-Headers"
	// HeaderMaxAge Indicates the number of seconds the information provided can be cached.
	// See: Fetch Standard § 3.2.3. HTTP responses.
	HeaderMaxAge = "Access-Control-Max-Age"
	// HeaderRequestMethod indicates which method a future CORS request to the same resource might use.
	// See: Fetch Standard § 3.2.2. HTTP requests.
	HeaderRequestMethod = "Access-Control-Request-Method"
	// HeaderRequestHeaders indicates which headers a future CORS request to the same resource might use.
	// See: Fetch Standard § 3.2.2. HTTP requests.
	HeaderRequestHeaders = "Access-Control-Request-Headers"

	// HeaderValueWildcard represents the wildcard CORS response, which allows any method,
	// header, or origin.
	// See: Fetch Standard § 3.2.4. HTTP new-header syntax.
	HeaderValueWildcard = "*"

	// DefaultMaxAge is the default time that a client should cache a
	// CORS preflight response.
	// See: Fetch Standard § 3.2.3. HTTP responses.
	DefaultMaxAge = 5
)

// Request represents a CORS request, which may or may not be a preflight request.
type Request http.Request

// IsPreflight determines if a request is a CORS preflight request.
// See: Fetch Standard § 3.2.2. HTTP requests.
func (r *Request) IsPreflight() bool {
	return r.Method == http.MethodOptions &&
		r.Header.Get(HeaderOrigin) != "" &&
		r.Header.Get(HeaderRequestMethod) != "" &&
		r.Header.Get(HeaderRequestHeaders) != ""
}

// Options represents the potential CORS options a server can return to its clients.
type Options struct {
	AllowCredentials bool
	AllowHeaders     []string
	AllowMethods     []string
	AllowOrigins     []string
	ExposeHeaders    []string
	MaxAge           int

	cache map[string]string
}

// NewOptions returns a properly initialized Options pointer.
func NewOptions() *Options {
	return &Options{
		AllowCredentials: false,
		AllowHeaders:     []string{},
		AllowMethods:     []string{},
		AllowOrigins:     []string{},
		ExposeHeaders:    []string{},
		MaxAge:           DefaultMaxAge,

		cache: nil,
	}
}

// GetAllowOrigin returns the appropriate Access-Control-Allow-Origin header.
// If the wildcard is present, it will be used instead of the request's
// Origin header. An empty string represents that no Access-Control-Allow-Origin
// should be returned to the client. The Access-Control-Allow-Origin header should
// be returned on all allowed CORS requests.
// See: Fetch Standard § 3.2.3. HTTP responses.
//
// If the client's credentials mode is "include", wildcard values will result in
// a client side failure.
// See: Fetch Standard § 3.2.5. CORS protocol and credentials.
func (o *Options) GetAllowOrigin(request *Request) string {
	origin := request.Header.Get(HeaderOrigin)
	result := ""

	for _, ao := range o.AllowOrigins {
		switch ao {
		case HeaderValueWildcard:
			return HeaderValueWildcard
		case origin:
			result = origin
		}
	}

	return result
}

// GetAllowCredentials returns the appropriate Access-Control-Allow-Credentials header.
// An empty string represents that no Access-Control-Allow-Credentials header should be
// returned to the client. The Access-Control-Allow-Credentials header should be returned on
// all CORS requests if credentials are allowed, regardless of the value of the client's
// credentials mode.
// See: Fetch Standard § 3.2.3. HTTP responses.
func (o *Options) GetAllowCredentials() string {
	if o.AllowCredentials {
		return "true"
	}

	return ""
}

// GetAllowMethods returns the appropriate Access-Control-Allow-Methods header.
// If the wildcard is present, it will be used instead of a comma-separated list.
// An empty string represents that no Access-Control-Allow-Methods header should be
// returned to the client. The Access-Control-Allow-Methods header should be returned
// on preflight requests.
// See: Fetch Standard § 3.2.3. HTTP responses.
//
// If the client's credentials mode is "include", wildcard values will result in
// a client side failure.
// See: Fetch Standard § 3.2.4. HTTP new-header syntax.
func (o *Options) GetAllowMethods() string {
	for _, am := range o.AllowMethods {
		if am == HeaderValueWildcard {
			return HeaderValueWildcard
		}
	}

	return strings.Join(o.AllowMethods, ", ")
}

// GetAllowHeaders returns the appropriate Access-Control-Allow-Headers header.
// If the wildcard is present, it will be used instead of a comma-separated list.
// An empty string represents that no Access-Control-Allow-Headers header should
// be returned to the client. The Access-Control-Allow-Headers header should be
// returned on preflight requests.
// See: Fetch Standard § 3.2.3. HTTP responses.
//
// If the client's credentials mode is "include", wildcard values will result in
// a client side failure.
// See: Fetch Standard § 3.2.4. HTTP new-header syntax.
func (o *Options) GetAllowHeaders() string {
	for _, ah := range o.AllowHeaders {
		if ah == HeaderValueWildcard {
			return HeaderValueWildcard
		}
	}

	return strings.Join(o.AllowHeaders, ", ")
}

// GetMaxAge returns the appropriate Access-Control-Max-Age header. An empty
// string represents that no Access-Control-Max-Age header should be returned
// to the client. The Access-Control-Max-Age header should be returned on
// preflight requests.
// See: Fetch Standard § 3.2.3. HTTP responses.
func (o *Options) GetMaxAge() string {
	return strconv.Itoa(o.MaxAge)
}

// GetExposeHeaders returns the appropriate Access-Control-Expose-Headers header.
// If the wildcard is present, it will be used instead of a comma-separated list.
// An empty string represents that no Access-Control-Expose-Headers header should
// be returned to the client. The Access-Control-Expose-Headers header should be
// returned on CORS requests that are not preflight requests.
// See: Fetch Standard § 3.2.3. HTTP responses.
//
// If the client's credentials mode is "include", wildcard values will result in
// a client side failure.
// See: Fetch Standard § 3.2.4. HTTP new-header syntax.
func (o *Options) GetExposeHeaders() string {
	for _, em := range o.ExposeHeaders {
		if em == HeaderValueWildcard {
			return HeaderValueWildcard
		}
	}

	return strings.Join(o.ExposeHeaders, ", ")
}

// GetVary returns the appropriate Vary header. An empty string represents that
// the Vary header should not be modified. The Vary header should include Origin
// if the server has multiple allowed origins, unless the server uses the
// wildcard origin.
// See: Fetch Standard § CORS protocol and HTTP caches.
func (o *Options) GetVary() string {
	if len(o.AllowOrigins) > 1 {
		return HeaderOrigin
	}

	return ""
}

// NewHandler returns a http.Handler that can process CORS requests from the
// provided Options.
func (o *Options) NewHandler() http.Handler {
	o.cache = make(map[string]string)

	o.cache[HeaderAllowMethods] = o.GetAllowMethods()
	o.cache[HeaderAllowHeaders] = o.GetAllowHeaders()
	o.cache[HeaderExposeHeaders] = o.GetExposeHeaders()
	o.cache[HeaderMaxAge] = o.GetMaxAge()

	return (*handler)(o)
}

type handler Options

// ServeHTTP implements http.Handler for Options.
func (h *handler) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	o := (*Options)(h)
	r := (*Request)(req)

	if v := o.GetVary(); v != "" {
		rw.Header().Add(HeaderVary, v)
	}

	if v := o.GetAllowOrigin(r); v != "" {
		rw.Header().Set(HeaderAllowOrigin, v)
	}

	if v := o.GetAllowCredentials(); v != "" {
		rw.Header().Set(HeaderAllowCredentials, v)
	}

	if r.IsPreflight() {
		if v := o.cache[HeaderAllowMethods]; v != "" {
			rw.Header().Set(HeaderAllowMethods, v)
		}

		if v := o.cache[HeaderAllowHeaders]; v != "" {
			rw.Header().Set(HeaderAllowHeaders, v)
		}

		if v := o.cache[HeaderMaxAge]; v != "" {
			rw.Header().Set(HeaderMaxAge, v)
		}

		rw.WriteHeader(http.StatusNoContent)

		return
	}

	if v := o.cache[HeaderExposeHeaders]; v != "" {
		rw.Header().Set(HeaderExposeHeaders, v)
	}
}
