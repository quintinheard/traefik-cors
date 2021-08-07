// Package traefik exports a Traefik plugin for use with Traefik Pilot.
package traefik

import (
	"context"
	"net/http"

	"github.com/quintinheard/traefik-cors/cors"
)

// Config represents the plugin configuration.
type Config struct {
	AllowCredentials bool     `json:"allowCredentials,omitempty"`
	AllowHeaders     []string `json:"allowHeaders,omitempty"`
	AllowMethods     []string `json:"allowMethods,omitempty"`
	AllowOrigins     []string `json:"allowOrigins,omitempty"`
	ExposeHeaders    []string `json:"exposeHeaders,omitempty"`
	MaxAge           int      `json:"maxAge,omitempty"`
}

// CreateConfig creates the default plugin configuration.
func CreateConfig() *Config {
	return &Config{
		AllowCredentials: false,
		AllowHeaders:     []string{},
		AllowMethods:     []string{http.MethodHead, http.MethodGet, http.MethodPost},
		AllowOrigins:     []string{"*"},
		ExposeHeaders:    []string{},
		MaxAge:           cors.DefaultMaxAge,
	}
}

// CorsPlugin a Traefik plugin.
type CorsPlugin struct {
	next http.Handler
	name string
	cors http.Handler
}

// New create a new CORS plugin.
func New(ctx context.Context, next http.Handler, config *Config, name string) (http.Handler, error) {
	c := &cors.Options{
		AllowCredentials: config.AllowCredentials,
		AllowHeaders:     config.AllowHeaders,
		AllowMethods:     config.AllowMethods,
		AllowOrigins:     config.AllowOrigins,
		ExposeHeaders:    config.ExposeHeaders,
		MaxAge:           config.MaxAge,
	}

	return &CorsPlugin{
		next: next,
		name: name,
		cors: c.NewHandler(),
	}, nil
}

func (c *CorsPlugin) ServeHTTP(rw http.ResponseWriter, req *http.Request) {
	c.cors.ServeHTTP(rw, req)

	if (*cors.Request)(req).IsPreflight() {
		return
	}

	c.next.ServeHTTP(rw, req)
}
