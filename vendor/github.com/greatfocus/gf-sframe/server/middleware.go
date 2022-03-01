package server

import (
	"context"
	"net"
	"net/http"
	"strings"
	"time"
)

var Limiter = NewThrottle()

// SetHeaders // prepare header response
func SetHeaders() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			(w).Header().Set("Content-Type", "application/json")
			(w).Header().Set("Access-Control-Allow-Methods", "POST, GET, OPTIONS, PUT, DELETE")
			(w).Header().Set("Access-Control-Allow-Headers", "Accept, Content-Type, Content-Length, Accept-Encoding, X-CSRF-JWT, Authorization, request-id")

			// continue
			h.ServeHTTP(w, r)
		})
	}
}

// Preflight validates request for jwt header
func Preflight() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			if (*r).Method == "OPTIONS" {
				(w).WriteHeader(http.StatusOK)
				return
			}
			// continue
			h.ServeHTTP(w, r)
		})
	}
}

// IsAllowedOrigin enable cors within the http handler
func IsAllowedOrigin(allowedOrigin string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			allowed := false
			origin := ip(r)

			// check if cors is available in list
			origins := strings.Split(allowedOrigin, ",")
			for _, v := range origins {
				if v == origin {
					allowed = true
				}
			}

			// allow cors if found
			if !allowed {
				(w).WriteHeader(http.StatusForbidden)
				return
			}

			// continue
			h.ServeHTTP(w, r)
		})
	}
}

// IsThrottle handle limits and rates
func IsThrottle() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// limit us requests per second
			if Limiter.IsThrottled(ip(r)) {
				(w).WriteHeader(http.StatusTooManyRequests)
				return
			}

			// continue
			h.ServeHTTP(w, r)
		})
	}
}

// IsAllowedIPs allow specific IP address
func IsAllowedIPs(allowedIps string) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			allowed := false
			// check if ip is available in list
			ips := strings.Split(allowedIps, ",")
			for _, v := range ips {
				if v == ip(r) {
					allowed = true
				}
			}
			// allow ip if found
			if !allowed {
				(w).WriteHeader(http.StatusForbidden)
				return
			}

			// continue
			h.ServeHTTP(w, r)
		})
	}
}

// IsAuthorized validate if users is allowed to access route
func IsAuthorized(jwt JWT) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			token, err := jwt.GetTokenInfo(r)
			if err != nil {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			var allowed bool
			var pattern = r.URL.Path
			for _, value := range token.Permissions {
				if value == pattern {
					allowed = true
				}
			}

			if !allowed {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// continue
			h.ServeHTTP(w, r)
		})
	}
}

// ProcessTimeout put a time limit for the handler process duration and will give an error response if timeout
func ProcessTimeout(timeout time.Duration) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			ctx, cancel := context.WithTimeout(r.Context(), timeout)
			defer cancel()

			r = r.WithContext(ctx)

			processDone := make(chan bool)
			go func() {
				h.ServeHTTP(w, r)
				processDone <- true
			}()

			select {
			case <-ctx.Done():
				w.WriteHeader(http.StatusRequestTimeout)
				return
			case <-processDone:

			}
		})
	}
}

// IsAuthenticated validates request for jwt header
func IsAuthenticated(jwt JWT) Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {

			// validate jwt
			if !jwt.IsValidToken(r) {
				w.WriteHeader(http.StatusUnauthorized)
				return
			}

			// continue
			h.ServeHTTP(w, r)
		})
	}
}

// NoAuthentication access without authentications
func NoAuthentication() Middleware {
	return func(h http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			// continue
			h.ServeHTTP(w, r)
		})
	}
}

// Middleware strct
type Middleware func(http.Handler) http.Handler

// Use middleware
func Use(h http.Handler, m ...Middleware) http.Handler {
	if len(m) < 1 {
		return h
	}
	wrapped := h
	for i := len(m) - 1; i >= 0; i-- {
		wrapped = m[i](wrapped)
	}
	return wrapped
}

func ip(r *http.Request) string {
	ip, _, _ := net.SplitHostPort(r.RemoteAddr)
	return ip
}
