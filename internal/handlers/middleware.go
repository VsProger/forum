package handlers

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"strings"
	"sync"
	"time"

	"github.com/VsProger/snippetbox/logger"
)

var logg = logger.NewLogger()

const (
	requestsPerMinute = 60
	windowSize        = time.Minute
	cleanupInterval   = time.Minute * 5
	clientTimeout     = time.Minute * 10
)

type client struct {
	requests    int
	windowStart time.Time
}

type RateLimiter struct {
	clients map[string]*client
	mu      sync.Mutex
}

func NewRateLimiter() *RateLimiter {
	rl := &RateLimiter{
		clients: make(map[string]*client),
	}
	go rl.cleanupExpiredClients()
	return rl
}

func (rl *RateLimiter) Allow(ip string) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()

	now := time.Now()
	if c, exists := rl.clients[ip]; exists {
		if now.Sub(c.windowStart) > windowSize {
			c.requests = 1
			c.windowStart = now
			return true
		}

		if c.requests < requestsPerMinute {
			c.requests++
			return true
		}

		return false
	}

	rl.clients[ip] = &client{
		requests:    1,
		windowStart: now,
	}
	return true
}

func (rl *RateLimiter) cleanupExpiredClients() {
	for {
		time.Sleep(cleanupInterval)
		rl.mu.Lock()
		now := time.Now()
		for ip, c := range rl.clients {
			if now.Sub(c.windowStart) > clientTimeout {
				delete(rl.clients, ip)
			}
		}
		rl.mu.Unlock()
	}
}

func (rl *RateLimiter) Middleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		if !rl.Allow(ip) {
			http.Error(w, "Too Many Requests", http.StatusTooManyRequests)
			return
		}
		next.ServeHTTP(w, r)
	})
}

func getIP(r *http.Request) string {
	xff := r.Header.Get("X-Forwarded-For")
	if xff != "" {
		ips := strings.Split(xff, ",")
		return strings.TrimSpace(ips[0])
	}

	xri := r.Header.Get("X-Real-IP")
	if xri != "" {
		return xri
	}

	ip, _, err := net.SplitHostPort(r.RemoteAddr)
	if err != nil {
		return r.RemoteAddr
	}
	return ip
}

func (h *Handler) LoggingMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		ip := getIP(r)
		method := r.Method
		uri := r.URL.RequestURI()
		proto := r.Proto
		start := time.Now()

		lrw := NewLoggingResponseWriter(w)

		next.ServeHTTP(lrw, r)

		duration := time.Since(start)

		logg.Info(fmt.Sprintf("%s %s %s %s %d %d %v",
			ip,
			method,
			uri,
			proto,
			lrw.statusCode,
			lrw.responseSize,
			duration,
		))
	})
}

type LoggingResponseWriter struct {
	http.ResponseWriter
	statusCode   int
	responseSize int
}

func NewLoggingResponseWriter(w http.ResponseWriter) *LoggingResponseWriter {
	return &LoggingResponseWriter{w, http.StatusOK, 0}
}

func (lrw *LoggingResponseWriter) WriteHeader(code int) {
	lrw.statusCode = code
	lrw.ResponseWriter.WriteHeader(code)
}

func (lrw *LoggingResponseWriter) Write(b []byte) (int, error) {
	size, err := lrw.ResponseWriter.Write(b)
	lrw.responseSize += size
	return size, err
}

func (h *Handler) AllHandler(next http.Handler) http.Handler {
	rateLimiter := NewRateLimiter()

	handler := rateLimiter.Middleware(next)

	handler = h.LoggingMiddleware(handler)
	handler = secureHeaders(handler)

	return handler
}

func (h *Handler) AuthMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		sessionCookie, err := r.Cookie("session")
		if err != nil || sessionCookie.Value == "" {

			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		user, err := h.service.GetUserByToken(sessionCookie.Value)
		if err != nil {

			http.Redirect(w, r, "/login", http.StatusFound)
			return
		}

		ctx := r.Context()
		ctx = context.WithValue(ctx, "user", user)

		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func (h *Handler) RoleMiddleware(requiredRoles []string, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		session, err := r.Cookie("session")
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		user, err := h.service.GetUserByToken(session.Value)
		if err != nil {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}

		hasAccess := false
		for _, role := range requiredRoles {
			if user.Role == role {
				hasAccess = true
				break
			}
		}

		if !hasAccess {
			http.Error(w, "Forbidden", http.StatusForbidden)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func secureHeaders(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		//w.Header().Set("Content-Security-Policy", "default-src 'self'; style-src 'self' fonts.googleapis.com; font-src fonts.gstatic.com")
		w.Header().Set("Referrer-Policy", "origin-when-cross-origin")
		w.Header().Set("X-Content-Type-Options", "nosniff")
		w.Header().Set("X-Frame-Options", "deny")
		w.Header().Set("X-XSS-Protection", "0")

		next.ServeHTTP(w, r)
	})
}
