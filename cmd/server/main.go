package main

import (
	"context"
	"errors"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"movieapi/api"
	"movieapi/internal/config"
	"movieapi/internal/movie"
	"movieapi/internal/service"
	grpcTransport "movieapi/internal/transport/grpc"
	httpTransport "movieapi/internal/transport/http"
	"google.golang.org/grpc"
)

func main() {
	cfg := config.Load()
	svc := service.New(movie.NewClient(cfg.MovieBaseURL))
	httpServer := &http.Server{Addr: cfg.HTTPAddr, Handler: httpTransport.NewRouter(svc), ReadHeaderTimeout: 5 * time.Second}
	grpcServer := grpc.NewServer()
	api.RegisterMovieServiceServer(grpcServer, grpcTransport.NewServer(svc))
	grpcAddr := os.Getenv("GRPC_ADDR")
	if grpcAddr == "" { grpcAddr = ":9090" }
	listener, err := net.Listen("tcp", grpcAddr)
	if err != nil { log.Fatal(err) }
	go func() { if err := grpcServer.Serve(listener); err != nil { log.Printf("gRPC server: %v", err) } }()
	go func() { if err := httpServer.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) { log.Printf("HTTP server: %v", err) } }()
	log.Printf("HTTP listening on %s; gRPC listening on %s", cfg.HTTPAddr, grpcAddr)
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)
	<-stop
	ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
	defer cancel()
	_ = httpServer.Shutdown(ctx)
	grpcServer.GracefulStop()
}
