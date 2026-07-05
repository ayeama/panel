package middleware

import (
	"bufio"
	"fmt"
	"net"
	"net/http"
	"time"
)

type statusRecorder struct {
	http.ResponseWriter
	status int
	wrote  bool
}

func (r *statusRecorder) WriteHeader(code int) {
	if !r.wrote {
		r.status = code
		r.wrote = true
	}
	r.ResponseWriter.WriteHeader(code)
}

func (r *statusRecorder) Write(b []byte) (int, error) {
	if !r.wrote {
		r.WriteHeader(http.StatusOK)
	}
	return r.ResponseWriter.Write(b)
}

func (r *statusRecorder) Hijack() (net.Conn, *bufio.ReadWriter, error) {
	hj, ok := r.ResponseWriter.(http.Hijacker)
	if !ok {
		return nil, nil, http.ErrNotSupported
	}
	return hj.Hijack()
}

func (r *statusRecorder) Flush() {
	if fl, ok := r.ResponseWriter.(http.Flusher); ok {
		fl.Flush()
	}
}

func Log(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		rec := &statusRecorder{ResponseWriter: w}

		timeStart := time.Now().UTC()
		next.ServeHTTP(rec, r)
		timeEnd := time.Now().UTC()

		fmt.Printf(
			"handled request host=\"%s\" method=\"%s\" path=\"%s\" query=\"%s\" status=%d time=%d\n",
			r.Host,
			r.Method,
			r.URL.Path,
			r.URL.RawQuery,
			rec.status,
			timeEnd.Sub(timeStart).Milliseconds(),
		)
	})
}
