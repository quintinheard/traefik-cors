package cors_test

import (
	"fmt"
	"net/http"

	"github.com/quintinheard/traefik-cors/cors"
)

func ExampleRequest() {
	o := cors.NewOptions()

	// an example http.HandlerFunc implementation
	_ = func(_ http.ResponseWriter, req *http.Request) {
		o.GetAllowOrigin((*cors.Request)(req))
	}
}

func ExampleRequest_IsPreflight() {
	// an example http.HandlerFunc implementation
	_ = func(rw http.ResponseWriter, req *http.Request) {
		corsReq := (*cors.Request)(req)
		if corsReq.IsPreflight() {
			rw.WriteHeader(http.StatusNoContent)
			return
		}
	}
}

func ExampleOptions() {
	o := cors.Options{
		AllowCredentials: false,
		AllowHeaders:     []string{},
		AllowMethods:     []string{http.MethodGet, http.MethodPost},
		AllowOrigins:     []string{},
		ExposeHeaders:    []string{},
		MaxAge:           cors.DefaultMaxAge,
	}

	fmt.Println(o.GetAllowMethods())
	// Output:
	// GET, POST
}

func ExampleNewOptions() {
	o := cors.NewOptions()

	fmt.Println(o.GetMaxAge())
	// Output:
	// 5
}

func ExampleOptions_GetAllowOrigin() {
	o := cors.NewOptions()

	// an example http.HandlerFunc implementation
	_ = func(rw http.ResponseWriter, req *http.Request) {
		if header := o.GetAllowOrigin((*cors.Request)(req)); header != "" {
			rw.Header().Set(cors.HeaderAllowOrigin, header)
		}
	}
}

func ExampleOptions_GetAllowCredentials() {
	o := cors.NewOptions()

	// an example http.HandlerFunc implementation
	_ = func(rw http.ResponseWriter, _ *http.Request) {
		if header := o.GetAllowCredentials(); header != "" {
			rw.Header().Set(cors.HeaderAllowCredentials, header)
		}
	}
}

func ExampleOptions_GetAllowMethods() {
	o := cors.NewOptions()

	// an example http.HandlerFunc implementation
	_ = func(rw http.ResponseWriter, _ *http.Request) {
		if header := o.GetAllowMethods(); header != "" {
			rw.Header().Set(cors.HeaderAllowMethods, header)
		}
	}
}

func ExampleOptions_GetAllowHeaders() {
	o := cors.NewOptions()

	// an example http.HandlerFunc implementation
	_ = func(rw http.ResponseWriter, _ *http.Request) {
		if header := o.GetAllowHeaders(); header != "" {
			rw.Header().Set(cors.HeaderAllowHeaders, header)
		}
	}
}

func ExampleOptions_GetMaxAge() {
	o := cors.NewOptions()

	// an example http.HandlerFunc implementation
	_ = func(rw http.ResponseWriter, _ *http.Request) {
		if header := o.GetMaxAge(); header != "" {
			rw.Header().Set(cors.HeaderMaxAge, header)
		}
	}
}

func ExampleOptions_GetExposeHeaders() {
	o := cors.NewOptions()

	// an example http.HandlerFunc implementation
	_ = func(rw http.ResponseWriter, _ *http.Request) {
		if header := o.GetExposeHeaders(); header != "" {
			rw.Header().Set(cors.HeaderExposeHeaders, header)
		}
	}
}

func ExampleOptions_GetVary() {
	o := cors.NewOptions()

	// an example http.HandlerFunc implementation
	_ = func(rw http.ResponseWriter, _ *http.Request) {
		if header := o.GetVary(); header != "" {
			rw.Header().Add(cors.HeaderVary, header)
		}
	}
}

func ExampleOptions_NewHandler() {
	h := cors.NewOptions().NewHandler()

	_ = http.ListenAndServe(":80", h)
}

func ExampleHandler_ServeHTTP() {
	h := cors.NewOptions().NewHandler()

	_ = http.ListenAndServe(":80", h)
}
