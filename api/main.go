package main

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	functionDeployments = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "function_deployments_total",
			Help: "Total number of function deployments",
		},
		[]string{"function", "status"},
	)
	functionInvocations = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "function_invocations_total",
			Help: "Total number of function invocations",
		},
		[]string{"function"},
	)
	functionDuration = prometheus.NewHistogramVec(
		prometheus.HistogramOpts{
			Name:    "function_duration_seconds",
			Help:    "Function execution duration in seconds",
			Buckets: prometheus.DefBuckets,
		},
		[]string{"function"},
	)
	coldStarts = prometheus.NewCounterVec(
		prometheus.CounterOpts{
			Name: "function_cold_starts_total",
			Help: "Total number of cold starts",
		},
		[]string{"function"},
	)
)

func init() {
	prometheus.MustRegister(functionDeployments)
	prometheus.MustRegister(functionInvocations)
	prometheus.MustRegister(functionDuration)
	prometheus.MustRegister(coldStarts)
}

type Server struct {
	k8sClient *KubernetesClient
	port      string
}

func NewServer(port string) (*Server, error) {
	k8sClient, err := NewKubernetesClient()
	if err != nil {
		return nil, fmt.Errorf("failed to create Kubernetes client: %w", err)
	}

	return &Server{
		k8sClient: k8sClient,
		port:      port,
	}, nil
}

func (s *Server) Start() error {
	r := mux.NewRouter()

	// Health endpoints
	r.HandleFunc("/health", s.healthHandler).Methods("GET")
	r.HandleFunc("/ready", s.readyHandler).Methods("GET")

	// Function management
	r.HandleFunc("/api/v1/functions", s.listFunctionsHandler).Methods("GET")
	r.HandleFunc("/api/v1/functions", s.createFunctionHandler).Methods("POST")
	r.HandleFunc("/api/v1/functions/{name}", s.getFunctionHandler).Methods("GET")
	r.HandleFunc("/api/v1/functions/{name}", s.updateFunctionHandler).Methods("PUT")
	r.HandleFunc("/api/v1/functions/{name}", s.deleteFunctionHandler).Methods("DELETE")

	// Function invocation
	r.HandleFunc("/api/v1/functions/{name}/invoke", s.invokeFunctionHandler).Methods("POST")

	// Metrics
	r.HandleFunc("/api/v1/functions/{name}/metrics", s.functionMetricsHandler).Methods("GET")

	// CORS middleware
	r.Use(corsMiddleware)

	log.Printf("Starting API server on port %s", s.port)
	return http.ListenAndServe(":"+s.port, r)
}

func (s *Server) healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func (s *Server) readyHandler(w http.ResponseWriter, r *http.Request) {
	// Check if we can connect to Kubernetes
	if err := s.k8sClient.Ping(); err != nil {
		w.WriteHeader(http.StatusServiceUnavailable)
		json.NewEncoder(w).Encode(map[string]string{"status": "not ready", "error": err.Error()})
		return
	}
	json.NewEncoder(w).Encode(map[string]string{"status": "ready"})
}

func (s *Server) listFunctionsHandler(w http.ResponseWriter, r *http.Request) {
	functions, err := s.k8sClient.ListFunctions(r.Context())
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(functions)
}

func (s *Server) createFunctionHandler(w http.ResponseWriter, r *http.Request) {
	var function Function
	if err := json.NewDecoder(r.Body).Decode(&function); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if err := s.k8sClient.CreateFunction(r.Context(), &function); err != nil {
		functionDeployments.WithLabelValues(function.Name, "failed").Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	functionDeployments.WithLabelValues(function.Name, "success").Inc()

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)
	json.NewEncoder(w).Encode(function)
}

func (s *Server) getFunctionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	function, err := s.k8sClient.GetFunction(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusNotFound)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(function)
}

func (s *Server) updateFunctionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	var function Function
	if err := json.NewDecoder(r.Body).Decode(&function); err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	function.Name = name

	if err := s.k8sClient.UpdateFunction(r.Context(), &function); err != nil {
		functionDeployments.WithLabelValues(function.Name, "failed").Inc()
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	functionDeployments.WithLabelValues(function.Name, "success").Inc()

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(function)
}

func (s *Server) deleteFunctionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	if err := s.k8sClient.DeleteFunction(r.Context(), name); err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.WriteHeader(http.StatusNoContent)
}

func (s *Server) invokeFunctionHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	start := time.Now()
	functionInvocations.WithLabelValues(name).Inc()

	result, err := s.k8sClient.InvokeFunction(r.Context(), name, r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	duration := time.Since(start).Seconds()
	functionDuration.WithLabelValues(name).Observe(duration)

	w.Header().Set("Content-Type", "application/json")
	w.Write(result)
}

func (s *Server) functionMetricsHandler(w http.ResponseWriter, r *http.Request) {
	vars := mux.Vars(r)
	name := vars["name"]

	metrics, err := s.k8sClient.GetFunctionMetrics(r.Context(), name)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	w.Header().Set("Content-Type", "application/json")
	json.NewEncoder(w).Encode(metrics)
}

func corsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		w.Header().Set("Access-Control-Allow-Origin", "*")
		w.Header().Set("Access-Control-Allow-Methods", "GET, POST, PUT, DELETE, OPTIONS")
		w.Header().Set("Access-Control-Allow-Headers", "Content-Type, Authorization")

		if r.Method == "OPTIONS" {
			w.WriteHeader(http.StatusOK)
			return
		}

		next.ServeHTTP(w, r)
	})
}

func startMetricsServer(port string) {
	http.Handle("/metrics", promhttp.Handler())
	log.Printf("Starting metrics server on port %s", port)
	if err := http.ListenAndServe(":"+port, nil); err != nil {
		log.Fatalf("Failed to start metrics server: %v", err)
	}
}

func main() {
	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	metricsPort := os.Getenv("METRICS_PORT")
	if metricsPort == "" {
		metricsPort = "9090"
	}

	// Start metrics server in background
	go startMetricsServer(metricsPort)

	server, err := NewServer(port)
	if err != nil {
		log.Fatalf("Failed to create server: %v", err)
	}

	if err := server.Start(); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
