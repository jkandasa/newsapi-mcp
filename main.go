package main

import (
	"context"
	"flag"
	"fmt"
	"net/http"
	"os"
	"os/signal"
	"syscall"

	"github.com/mark3labs/mcp-go/server"
	"go.uber.org/zap"
)

var version = "devel" // overridden at build time via -ldflags "-X main.version=..."

func main() {
	var (
		showVersion = flag.Bool("version", false, "Print version and exit")
		transport   = flag.String("transport", "stdio", "Transport type: stdio or http")
		addr        = flag.String("addr", ":8080", "HTTP listen address (used with -transport=http)")
		certFile    = flag.String("cert", "", "TLS certificate file (enables HTTPS)")
		keyFile     = flag.String("key", "", "TLS key file (enables HTTPS)")
	)
	flag.Parse()

	if *showVersion {
		fmt.Println("newsapi-mcp", version)
		os.Exit(0)
	}

	log, _ := zap.NewProduction()
	defer log.Sync()

	apiKey := os.Getenv("NEWSAPI_KEY")
	if apiKey == "" {
		log.Fatal("NEWSAPI_KEY environment variable is required")
	}

	client := newClient(apiKey)

	s := server.NewMCPServer("newsapi-mcp", version,
		server.WithToolCapabilities(false),
	)
	registerTools(s, client)

	switch *transport {
	case "stdio":
		srv := server.NewStdioServer(s)
		if err := srv.Listen(context.Background(), os.Stdin, os.Stdout); err != nil {
			log.Fatal("stdio server error", zap.Error(err))
		}

	case "http":
		tls := *certFile != "" && *keyFile != ""
		if (*certFile == "") != (*keyFile == "") {
			log.Fatal("both -cert and -key must be provided together")
		}

		scheme := "http"
		if tls {
			scheme = "https"
		}
		baseURL := scheme + "://" + *addr
		sseSrv := server.NewSSEServer(s, server.WithBaseURL(baseURL))

		log.Info("newsapi-mcp started",
			zap.String("version", version),
			zap.String("transport", "http"),
			zap.Bool("tls", tls),
			zap.String("base_url", baseURL),
			zap.String("sse_endpoint", baseURL+"/sse"),
			zap.String("message_endpoint", baseURL+"/message"),
		)

		ctx, stop := signal.NotifyContext(context.Background(), os.Interrupt, syscall.SIGTERM)
		defer stop()

		httpSrv := &http.Server{Addr: *addr, Handler: sseSrv}
		errCh := make(chan error, 1)
		go func() {
			if tls {
				errCh <- httpSrv.ListenAndServeTLS(*certFile, *keyFile)
			} else {
				errCh <- httpSrv.ListenAndServe()
			}
		}()

		select {
		case err := <-errCh:
			log.Fatal("http server error", zap.Error(err))
		case <-ctx.Done():
			log.Info("shutting down")
			httpSrv.Shutdown(context.Background()) //nolint:errcheck
		}

	default:
		log.Fatal("unknown transport", zap.String("transport", *transport),
			zap.String("hint", "use 'stdio' or 'http'"))
	}
}
