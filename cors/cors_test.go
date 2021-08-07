package cors_test

import (
	"context"
	"net/http"
	"net/http/httptest"
	"testing"

	"github.com/quintinheard/traefik-cors/cors"
	"github.com/stretchr/testify/require"
)

func TestRequest_IsPreflight(t *testing.T) {
	req := (*cors.Request)(httptest.NewRequest(http.MethodOptions, "https://cors.example.com/api/", nil))
	req.Header.Set(cors.HeaderOrigin, "https://example.com")
	req.Header.Set(cors.HeaderRequestHeaders, "Content-Type")
	req.Header.Set(cors.HeaderRequestMethod, http.MethodGet)

	require.Equal(t, true, req.IsPreflight())
}

func TestRequest_IsNotPreflight(t *testing.T) {
	req := (*cors.Request)(httptest.NewRequest(http.MethodOptions, "https://cors.example.com/api/", nil))
	req.Header.Set(cors.HeaderOrigin, "https://example.com")

	require.Equal(t, false, req.IsPreflight())
}

func TestHandler_ServeHTTP(t *testing.T) {
	o := cors.NewOptions()
	o.AllowOrigins = []string{"https://example.com"}
	o.AllowHeaders = []string{"Content-Type", "Authorization"}
	o.AllowMethods = []string{http.MethodGet, http.MethodPost}

	req, err := http.NewRequestWithContext(context.Background(), http.MethodOptions, "https://cors.example.com", nil)
	require.Nil(t, err)
	req.Header.Set(cors.HeaderOrigin, "https://example.com")
	req.Header.Set(cors.HeaderRequestHeaders, "Content-Type")
	req.Header.Set(cors.HeaderRequestMethod, http.MethodGet)

	rec := httptest.NewRecorder()

	o.NewHandler().ServeHTTP(rec, req)

	res := rec.Result()
	require.Equal(t, http.StatusNoContent, res.StatusCode)
	require.Equal(t, "https://example.com", res.Header.Get(cors.HeaderAllowOrigin))
	require.Equal(t, "Content-Type, Authorization", res.Header.Get(cors.HeaderAllowHeaders))
	require.Equal(t, "GET, POST", res.Header.Get(cors.HeaderAllowMethods))
	require.Nil(t, res.Body.Close())
}
