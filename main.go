package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"

	"golang.org/x/sync/errgroup"
)

func main() {
	// err := http.ListenAndServe(
	// 	":18080",
	// 	http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
	// 		fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
	// 	}),
	// )

	// if err != nil {
	// 	fmt.Printf("Failed to terminate server: %v", err)
	// 	os.Exit(1)
	// }

	if len(os.Args) != 2 {
		log.Printf("need port number\n")
		os.Exit(1)
	}
	port := os.Args[1]
	listener, err := net.Listen("tcp", ":"+port)
	if err != nil {
		log.Fatalf("failed to listen port %s: %v", port, err)
	}
	if err := run(context.Background(), listener); err != nil {
		log.Printf("failed to terminate server: %v", err)
	}

	// s.ListenAndServe()
}

func run(ctx context.Context, listner net.Listener) error {
	s := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}

	eg, ctx := errgroup.WithContext(ctx)

	eg.Go(func() error {
		if err := s.Serve(listner); err != nil &&
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	<-ctx.Done()
	if err := s.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}

	return eg.Wait()
}
