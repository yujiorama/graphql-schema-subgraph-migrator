package main

import (
	"bytes"
	"fmt"
	"go.uber.org/zap"
	"io"
	"log"
	"net/http"
	"os"
	"time"

	"example.venus/graph"
	"github.com/99designs/gqlgen/graphql/handler"
	"github.com/99designs/gqlgen/graphql/handler/extension"
	"github.com/99designs/gqlgen/graphql/handler/lru"
	"github.com/99designs/gqlgen/graphql/handler/transport"
	"github.com/99designs/gqlgen/graphql/playground"
	"github.com/vektah/gqlparser/v2/ast"
)

const defaultPort = "4002"

type rwWrapper struct {
	rw http.ResponseWriter
	mw io.Writer
}

func newRwWrapper(rw http.ResponseWriter, buf io.Writer) *rwWrapper {
	return &rwWrapper{
		rw: rw,
		mw: io.MultiWriter(rw, buf),
	}
}

func (r *rwWrapper) Header() http.Header {
	return r.rw.Header()
}

func (r *rwWrapper) Write(i []byte) (int, error) {
	return r.mw.Write(i)
}

func (r *rwWrapper) WriteHeader(i int) {
	r.rw.WriteHeader(i)
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = defaultPort
	}

	srv := handler.New(graph.NewExecutableSchema(graph.Config{Resolvers: graph.NewResolver()}))

	srv.AddTransport(transport.Options{})
	srv.AddTransport(transport.GET{})
	srv.AddTransport(transport.POST{})

	srv.SetQueryCache(lru.New[*ast.QueryDocument](1000))

	srv.Use(extension.Introspection{})
	srv.Use(extension.AutomaticPersistedQuery{
		Cache: lru.New[string](100),
	})

	http.Handle("/", playground.Handler("GraphQL playground", "/query"))
	http.Handle("/query", srv)

	logger, _ := zap.NewDevelopment()
	defer logger.Sync()
	accessLogger := func(next http.Handler) http.Handler {
		return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			requestBody, _ := io.ReadAll(r.Body)
			r.Body = io.NopCloser(bytes.NewBuffer(requestBody))
			responseBody := &bytes.Buffer{}
			rww := newRwWrapper(w, responseBody)
			startTime := time.Now()
			next.ServeHTTP(rww, r)

			defer func() {
				logger.Sugar().Infow("access log",
					"http_method", r.Method,
					"path", r.URL,
					"remote_addr", r.RemoteAddr,
					"user_agent", r.UserAgent(),
					"request", string(requestBody),
					"response", string(responseBody.Bytes()),
					"response_duration_second", fmt.Sprintf("%.3f", time.Since(startTime).Seconds()),
				)
			}()
		})
	}

	log.Printf("connect to http://localhost:%s/ for GraphQL playground", port)
	log.Fatal(http.ListenAndServe(":"+port, accessLogger(srv)))
}
