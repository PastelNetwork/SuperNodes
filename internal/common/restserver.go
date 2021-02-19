package common

import (
	"context"
	"fmt"
	"net/http"
	"sync"
)

func home(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Hello from Snippetbox"))
}

// Add a showSnippet handler function.
func showSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Display a specific snippet..."))
}

// Add a createSnippet handler function.
func createSnippet(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Create a new snippet..."))
}
func StartRestServer(ctx context.Context, config *Config, logger *Logger, wg *sync.WaitGroup) func() error {

	mux := http.NewServeMux()
	mux.HandleFunc("/", home)
	mux.HandleFunc("/snippet", showSnippet)
	mux.HandleFunc("/snippet/create", createSnippet)

	return CreateServer("rest_server", ctx, config, logger, wg,
		//startServer
		func(ctx context.Context) error {
			return nil
		},
		//runServer
		func(ctx context.Context) error {
			if err := http.ListenAndServe(":4000", mux); err != http.ErrServerClosed {
				return fmt.Errorf("error starting Rest server: %w", err)
			}
			return nil
		},
		//stopServer
		func(ctx context.Context) error {
			//return server.Shutdown(ctx)
			return nil
		})
}

