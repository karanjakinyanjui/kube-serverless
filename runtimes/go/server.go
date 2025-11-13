package main

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"log"
	"net/http"
	"os"
	"plugin"
	"time"

	"github.com/gorilla/mux"
	"github.com/prometheus/client_golang/prometheus"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	coldStart = true
	handler   func(map[string]interface{}) (interface{}, error)

	invocations = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "function_invocations_total",
		Help: "Total function invocations",
	})
	duration = prometheus.NewHistogram(prometheus.HistogramOpts{
		Name: "function_duration_seconds",
		Help: "Function execution duration",
	})
	coldStarts = prometheus.NewCounter(prometheus.CounterOpts{
		Name: "function_cold_starts_total",
		Help: "Total cold starts",
	})
)

func init() {
	prometheus.MustRegister(invocations)
	prometheus.MustRegister(duration)
	prometheus.MustRegister(coldStarts)
}

type Event struct {
	Body    interface{}       `json:"body"`
	Headers map[string]string `json:"headers"`
	Method  string            `json:"method"`
	Path    string            `json:"path"`
	Query   map[string]string `json:"query"`
}

type Response struct {
	StatusCode int         `json:"statusCode"`
	Body       interface{} `json:"body"`
}

func loadFunction() {
	codePath := "/function/code"
	handlerName := os.Getenv("FUNCTION_HANDLER")
	if handlerName == "" {
		handlerName = "handler"
	}

	if _, err := os.Stat(codePath); err == nil {
		// Try to load as Go plugin
		p, err := plugin.Open(codePath + ".so")
		if err == nil {
			sym, err := p.Lookup(handlerName)
			if err == nil {
				handler = sym.(func(map[string]interface{}) (interface{}, error))
				log.Println("Function loaded successfully")
				coldStart = false
				return
			}
		}
	}

	log.Println("No function code found or failed to load, using echo handler")
	handler = func(event map[string]interface{}) (interface{}, error) {
		return map[string]interface{}{
			"statusCode": 200,
			"body":       event,
		}, nil
	}
}

func healthHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]string{"status": "healthy"})
}

func readyHandler(w http.ResponseWriter, r *http.Request) {
	json.NewEncoder(w).Encode(map[string]interface{}{
		"status":    "ready",
		"coldStart": coldStart,
	})
}

func invokeHandler(w http.ResponseWriter, r *http.Request) {
	startTime := time.Now()
	wasColdStart := coldStart

	if coldStart {
		coldStarts.Inc()
		coldStart = false
	}

	invocations.Inc()

	body, err := ioutil.ReadAll(r.Body)
	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	var bodyJSON interface{}
	if len(body) > 0 {
		json.Unmarshal(body, &bodyJSON)
	}

	headers := make(map[string]string)
	for k, v := range r.Header {
		if len(v) > 0 {
			headers[k] = v[0]
		}
	}

	query := make(map[string]string)
	for k, v := range r.URL.Query() {
		if len(v) > 0 {
			query[k] = v[0]
		}
	}

	event := map[string]interface{}{
		"body":    bodyJSON,
		"headers": headers,
		"method":  r.Method,
		"path":    r.URL.Path,
		"query":   query,
	}

	result, err := handler(event)
	if err != nil {
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	execDuration := time.Since(startTime).Seconds()
	duration.Observe(execDuration)

	w.Header().Set("Content-Type", "application/json")
	w.Header().Set("X-Function-Duration", fmt.Sprintf("%.3f", execDuration))
	w.Header().Set("X-Cold-Start", fmt.Sprintf("%t", wasColdStart))

	if resp, ok := result.(Response); ok {
		w.WriteHeader(resp.StatusCode)
		json.NewEncoder(w).Encode(resp.Body)
	} else if respMap, ok := result.(map[string]interface{}); ok {
		if statusCode, ok := respMap["statusCode"].(int); ok {
			w.WriteHeader(statusCode)
		}
		json.NewEncoder(w).Encode(result)
	} else {
		json.NewEncoder(w).Encode(result)
	}
}

func main() {
	loadFunction()

	r := mux.NewRouter()

	r.HandleFunc("/health", healthHandler).Methods("GET")
	r.HandleFunc("/ready", readyHandler).Methods("GET")
	r.HandleFunc("/", invokeHandler).Methods("POST")
	r.Handle("/metrics", promhttp.Handler()).Methods("GET")

	port := os.Getenv("PORT")
	if port == "" {
		port = "8080"
	}

	functionName := os.Getenv("FUNCTION_NAME")
	handlerName := os.Getenv("FUNCTION_HANDLER")

	log.Printf("Go runtime server listening on port %s", port)
	log.Printf("Function: %s", functionName)
	log.Printf("Handler: %s", handlerName)

	if err := http.ListenAndServe(":"+port, r); err != nil {
		log.Fatalf("Failed to start server: %v", err)
	}
}
