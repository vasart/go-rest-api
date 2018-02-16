package main

import (
	"log"
	"net/http"
	"time"
)

// responseLogger is wrapper of http.ResponseWriter that keeps track of its HTTP
// status code and body size
type responseLogger struct {
	w      http.ResponseWriter
	status int
	size   int
}

func (l *responseLogger) Header() http.Header {
	return l.w.Header()
}

func (l *responseLogger) Write(b []byte) (int, error) {
	size, err := l.w.Write(b)
	l.size += size
	return size, err
}

func (l *responseLogger) WriteHeader(s int) {
	l.w.WriteHeader(s)
	l.status = s
}

func (l *responseLogger) Status() int {
	return l.status
}

func (l *responseLogger) Size() int {
	return l.size
}

func (l *responseLogger) Flush() {
	f, ok := l.w.(http.Flusher)
	if ok {
		f.Flush()
	}
}

func recovery(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		// NOTE
		// here we trying to recover from panicking in nested handlers
		// e.g. getMgoSession() may panic and we can handle it here,
		// responding with Server Error to client
		// for details see https://blog.golang.org/defer-panic-and-recover
		defer func(w http.ResponseWriter) {
			if err := recover(); err != nil {
				log.Print("recovered after: ")
				log.Println(err)
				//handleServerError(w, http.StatusInternalServerError)
				return
			}
		}(w)

		next.ServeHTTP(w, r)
	})
}

func logging(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()

		addr := r.RemoteAddr
		if val := r.Header.Get("X-Real-IP"); len(val) > 0 {
			addr = val
		}

		lw := responseLogger{
			w:      w,
			status: 200,
			size:   0,
		}
		next.ServeHTTP(&lw, r)

		log.Printf("%s \"%s %s %s\" %d %d %s", addr, r.Method, r.RequestURI, r.Proto, lw.Status(), lw.Size(), time.Since(start))
	})
}
