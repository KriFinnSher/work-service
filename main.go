package main

import (
	"context"
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
	"google.golang.org/grpc"

	"github.com/demo/work-service/gen"
)

var (
	requestsTotal = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "work_service_requests_total",
			Help: "Total number of requests",
		},
		[]string{"method", "status"},
	)
	requestDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "work_service_request_duration_seconds",
			Help:    "Request duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"method"},
	)
	activeConnections = prometheus.NewGauge(prometheus.GaugeOpts{
		Name: "work_service_active_connections",
		Help: "Number of active connections",
	})
)

func init() {
	prometheus.MustRegister(requestsTotal, requestDuration, activeConnections)
}

type server struct {
	gen.UnimplementedDemoServiceServer
}

func (s *server) GetUser(ctx context.Context, req *gen.GetUserRequest) (*gen.GetUserResponse, error) {
	start := time.Now()
	defer func() {
		requestsTotal.WithLabelValues("GetUser", "success").Inc()
		requestDuration.WithLabelValues("GetUser").Observe(time.Since(start).Seconds())
	}()

	log.Printf("GetUser request: id=%s", req.Id)

	if req.Id == "" {
		requestsTotal.WithLabelValues("GetUser", "error").Inc()
		return nil, fmt.Errorf("user id is required")
	}

	return &gen.GetUserResponse{
		Id:    req.Id,
		Email: "user@example.com",
		Name:  "John Doe",
	}, nil
}

func (s *server) CreateUser(ctx context.Context, req *gen.CreateUserRequest) (*gen.CreateUserResponse, error) {
	start := time.Now()
	defer func() {
		requestsTotal.WithLabelValues("CreateUser", "success").Inc()
		requestDuration.WithLabelValues("CreateUser").Observe(time.Since(start).Seconds())
	}()

	log.Printf("CreateUser request: email=%s name=%s", req.Email, req.Name)

	if req.Email == "" || req.Name == "" {
		requestsTotal.WithLabelValues("CreateUser", "error").Inc()
		return nil, fmt.Errorf("email and name are required")
	}

	userID := fmt.Sprintf("user-%d", time.Now().Unix())
	return &gen.CreateUserResponse{Id: userID}, nil
}

func main() {
	log.Println("Starting work-service...")

	grpcPort := os.Getenv("GRPC_PORT")
	if grpcPort == "" {
		grpcPort = "8082"
	}
	httpPort := os.Getenv("HTTP_PORT")
	if httpPort == "" {
		httpPort = "8083"
	}

	go func() {
		mux := http.NewServeMux()
		mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[WARN] healthz endpoint called from %s", r.RemoteAddr)
			w.Write([]byte("ok v3"))
		})
		mux.HandleFunc("/test", func(w http.ResponseWriter, r *http.Request) {
			log.Printf("[WARN] test endpoint called from %s", r.RemoteAddr)
			_, _ = (&server{}).GetUser(context.Background(), &gen.GetUserRequest{
				Id: "123",
			})
			w.Write([]byte("ok"))
		})
		mux.Handle("/metrics", promhttp.Handler())
		log.Printf("HTTP server listening on :%s", httpPort)
		if err := http.ListenAndServe(":"+httpPort, mux); err != nil {
			log.Printf("HTTP server error: %v", err)
		}
	}()

	lis, err := net.Listen("tcp", ":"+grpcPort)
	if err != nil {
		log.Fatalf("failed to listen: %v", err)
	}

	s := grpc.NewServer()
	gen.RegisterDemoServiceServer(s, &server{})

	log.Printf("gRPC server listening on :%s", grpcPort)

	go func() {
		sigCh := make(chan os.Signal, 1)
		signal.Notify(sigCh, syscall.SIGINT, syscall.SIGTERM)
		<-sigCh
		log.Println("Shutting down gracefully...")
		s.GracefulStop()
	}()

	if err := s.Serve(lis); err != nil {
		log.Fatalf("failed to serve: %v", err)
	}
}
