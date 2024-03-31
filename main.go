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
}

func run(ctx context.Context, listner net.Listener) error {
	// net/http 패키지에서 *http.Server 타입을 사용하여 HTTP 서버 실행 및 종료
	server := &http.Server{
		Handler: http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
			fmt.Fprintf(w, "Hello, %s!", r.URL.Path[1:])
		}),
	}

	eg, ctx := errgroup.WithContext(ctx)

	// 새로운 고루틴에서 HTTP 서버를 실행
	eg.Go(func() error {
		if err := server.Serve(listner); err != nil &&
			// http.ErrServerClosed는 에러가 아닌 정상 종료된 것
			err != http.ErrServerClosed {
			log.Printf("failed to close: %+v", err)
			return err
		}
		return nil
	})

	// 채널로부터의 알림(종료 알림)을 기다림
	<-ctx.Done()
	if err := server.Shutdown(context.Background()); err != nil {
		log.Printf("failed to shutdown: %+v", err)
	}
	// Go 메서드로 실행한 다른 고루틴의 종료를 기다림
	return eg.Wait()
}
