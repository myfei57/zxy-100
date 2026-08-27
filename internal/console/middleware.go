package console

import (
	"fmt"
	"net/http"
	"time"
)

func (s *Server) withLogging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		elapsed := time.Since(start).Round(time.Millisecond)
		s.audit.Append("console", "request", fmt.Sprintf("%s %s %s", r.Method, r.URL.Path, elapsed.String()))
	})
}
