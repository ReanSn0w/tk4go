package web

import (
	"bytes"
	"net/http"
	"net/http/httputil"
	"strings"

	"github.com/ReanSn0w/tk4go/pkg/tools"
)

// OptionalMiddlewares returns a middleware that applies the given middlewares
func OptionalMiddlewares(active bool, middlewares ...func(http.Handler) http.Handler) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		if !active {
			return next
		}

		for i := len(middlewares) - 1; i >= 0; i-- {
			next = middlewares[i](next)
		}

		return next
	}
}

// LoggerMiddleware is a middleware that logs requests and responses
// middleware uses [DEBUG] level
func LoggerMiddleware(logger tools.Logger) func(http.Handler) http.Handler {
	return func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			dumpBody := strings.Contains(r.Header.Get("Content-Type"), "json")
			reqDump, _ := httputil.DumpRequest(r, dumpBody)

			lrw := &loggedResponseWriter{
				statusCode: http.StatusOK,
				header:     make(http.Header),
				body:       new(bytes.Buffer),
			}

			next.ServeHTTP(lrw, r)
			lrw.rewriteTo(w)

			go func() {
				logger.Logf(
					"[DEBUG] http server dump\n\n %v\n---\n%v",
					string(reqDump), string(lrw.dump()))
			}()
		})
	}
}

type loggedResponseWriter struct {
	statusCode int
	header     http.Header
	body       *bytes.Buffer
}

func (l *loggedResponseWriter) Header() http.Header {
	return l.header
}

func (l *loggedResponseWriter) Write(b []byte) (int, error) {
	return l.body.Write(b)
}

func (l *loggedResponseWriter) WriteHeader(code int) {
	l.statusCode = code
}

func (l *loggedResponseWriter) rewriteTo(wr http.ResponseWriter) {
	wr.WriteHeader(l.statusCode)

	for k, v := range l.header {
		wr.Header().Set(k, v[0])
	}

	l.body.WriteTo(wr)
}

func (l *loggedResponseWriter) dump() []byte {
	buf := new(bytes.Buffer)

	buf.WriteString("HTTP/1.1 ")
	buf.WriteString(http.StatusText(l.statusCode))
	buf.WriteString("\n")

	for k, v := range l.header {
		buf.WriteString(k)
		buf.WriteString(": ")
		buf.WriteString(v[0])
		buf.WriteString("\n")
	}

	buf.WriteString("\n")
	l.header.Write(buf)
	buf.WriteString("\n")
	buf.Write(l.body.Bytes())

	return buf.Bytes()
}
