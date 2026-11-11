package main

import (
	"context"
	_ "embed"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"strconv"
	"strings"
	"sync"
	"sync/atomic"
	"syscall"
	"time"
)

//go:embed dashboard.html
var dashboardHTML []byte

var (
	version       = "2.0.0"
	startTime     = time.Now()
	requestCount  uint64
	errorCount    uint64
	cloudProvider = getEnv("CLOUD_PROVIDER", "unknown")
	clusterName   = getEnv("CLUSTER_NAME", "local-dev")
	environment   = getEnv("ENVIRONMENT", "development")
	podName       = getEnv("POD_NAME", "localhost")
	nodeName      = getEnv("NODE_NAME", "localhost")
)

type ResponseWriterInterceptor struct {
	http.ResponseWriter
	statusCode int
}

func (w *ResponseWriterInterceptor) WriteHeader(code int) {
	w.statusCode = code
	w.ResponseWriter.WriteHeader(code)
}

type MetricsCollector struct {
	mu           sync.RWMutex
	pathCounts   map[string]uint64
	statusCounts map[int]uint64
	latencies    map[string]time.Duration
}

var globalMetrics = &MetricsCollector{
	pathCounts:   make(map[string]uint64),
	statusCounts: make(map[int]uint64),
	latencies:    make(map[string]time.Duration),
}

func (m *MetricsCollector) Record(path string, status int, duration time.Duration) {
	m.mu.Lock()
	defer m.mu.Unlock()
	m.pathCounts[path]++
	m.statusCounts[status]++
	m.latencies[path] += duration
}

func loggingAndMetricsMiddleware(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		atomic.AddUint64(&requestCount, 1)

		wi := &ResponseWriterInterceptor{ResponseWriter: w, statusCode: http.StatusOK}
		next.ServeHTTP(wi, r)

		duration := time.Since(start)
		if wi.statusCode >= 400 {
			atomic.AddUint64(&errorCount, 1)
		}

		globalMetrics.Record(r.URL.Path, wi.statusCode, duration)

		log.Printf(`{"time":"%s","method":"%s","path":"%s","status":%d,"duration_ms":%.2f,"remote_addr":"%s"}`,
			time.Now().UTC().Format(time.RFC3339),
			r.Method,
			r.URL.Path,
			wi.statusCode,
			float64(duration.Microseconds())/1000.0,
			r.RemoteAddr,
		)
	})
}

func healthzHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	uptime := time.Since(startTime).Seconds()
	resp := map[string]interface{}{
		"status":         "UP",
		"uptime_seconds": uptime,
		"timestamp":      time.Now().UTC().Format(time.RFC3339),
		"version":        version,
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func readinessHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	resp := map[string]interface{}{
		"status":    "READY",
		"timestamp": time.Now().UTC().Format(time.RFC3339),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func detectSecurityPosture() (string, string) {
	p := strings.ToLower(cloudProvider)
	secretMounted := false
	if _, err := os.Stat("/mnt/secrets-store/DB-PASSWORD"); err == nil {
		secretMounted = true
	} else if _, err := os.Stat("/mnt/secrets-store"); err == nil {
		secretMounted = true
	}

	if strings.Contains(p, "aws") || os.Getenv("AWS_WEB_IDENTITY_TOKEN_FILE") != "" {
		return "IRSA (OIDC WebIdentity)", "KMS Envelope + Secrets CSI"
	}
	if strings.Contains(p, "azure") {
		if secretMounted {
			return "Entra Workload Identity", "Key Vault CSI Mounted (/mnt/secrets-store)"
		}
		return "Entra Managed Identity", "Azure Key Vault CSI Ready"
	}
	return "K8s ServiceAccount Token", "Zero-Trust Distroless Runtime"
}

func infoHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	idMode, secStatus := detectSecurityPosture()

	resp := map[string]interface{}{
		"service":        "multi-cloud-core-service",
		"version":        version,
		"cloud_provider": cloudProvider,
		"cluster_name":   clusterName,
		"environment":    environment,
		"pod_name":       podName,
		"node_name":      nodeName,
		"identity_mode":  idMode,
		"secret_status":  secStatus,
		"total_requests": atomic.LoadUint64(&requestCount),
		"runtime": map[string]interface{}{
			"go_version": runtime.Version(),
			"num_cpu":    runtime.NumCPU(),
			"goroutines": runtime.NumGoroutine(),
			"alloc_mb":   float64(memStats.Alloc) / (1024 * 1024),
			"sys_mb":     float64(memStats.Sys) / (1024 * 1024),
		},
		"uptime_seconds": time.Since(startTime).Seconds(),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func dataHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	items := []map[string]interface{}{
		{
			"id":       "svc-aws-01",
			"provider": "AWS",
			"tier":     "EKS v1.30 Managed NodeGroup (t3.medium Spot)",
			"status":   "ACTIVE",
		},
		{
			"id":       "svc-azure-02",
			"provider": "Azure",
			"tier":     "AKS v1.34 Azure CNI (Standard_D2s_v4 + Spot)",
			"status":   "ACTIVE",
		},
		{
			"id":       "svc-gcp-03",
			"provider": "GCP",
			"tier":     "GKE v1.29 Standard Cluster (IaC Ready / Standby)",
			"status":   "STANDBY",
		},
	}

	resp := map[string]interface{}{
		"timestamp": time.Now().UTC().Format(time.RFC3339),
		"count":     len(items),
		"data":      items,
		"requestId": fmt.Sprintf("req-%d", atomic.LoadUint64(&requestCount)),
	}
	_ = json.NewEncoder(w).Encode(resp)
}

func metricsHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "text/plain; version=0.0.4")

	var memStats runtime.MemStats
	runtime.ReadMemStats(&memStats)

	totalRequests := atomic.LoadUint64(&requestCount)
	totalErrors := atomic.LoadUint64(&errorCount)
	uptime := time.Since(startTime).Seconds()

	fmt.Fprintf(w, "# HELP http_requests_total Total number of HTTP requests processed\n")
	fmt.Fprintf(w, "# TYPE http_requests_total counter\n")
	fmt.Fprintf(w, "http_requests_total{service=\"core-service\",cloud=\"%s\",cluster=\"%s\"} %d\n",
		cloudProvider, clusterName, totalRequests)

	fmt.Fprintf(w, "# HELP http_errors_total Total number of HTTP requests returning error status (>=400)\n")
	fmt.Fprintf(w, "# TYPE http_errors_total counter\n")
	fmt.Fprintf(w, "http_errors_total{service=\"core-service\",cloud=\"%s\"} %d\n", cloudProvider, totalErrors)

	fmt.Fprintf(w, "# HELP app_uptime_seconds Process uptime in seconds\n")
	fmt.Fprintf(w, "# TYPE app_uptime_seconds gauge\n")
	fmt.Fprintf(w, "app_uptime_seconds{service=\"core-service\"} %.2f\n", uptime)

	fmt.Fprintf(w, "# HELP go_goroutines Number of goroutines currently existing\n")
	fmt.Fprintf(w, "# TYPE go_goroutines gauge\n")
	fmt.Fprintf(w, "go_goroutines %d\n", runtime.NumGoroutine())

	fmt.Fprintf(w, "# HELP go_memstats_alloc_bytes Bytes allocated and not yet freed\n")
	fmt.Fprintf(w, "# TYPE go_memstats_alloc_bytes gauge\n")
	fmt.Fprintf(w, "go_memstats_alloc_bytes %d\n", memStats.Alloc)

	globalMetrics.mu.RLock()
	defer globalMetrics.mu.RUnlock()

	fmt.Fprintf(w, "# HELP http_requests_by_path_total Total requests broken down by HTTP path\n")
	fmt.Fprintf(w, "# TYPE http_requests_by_path_total counter\n")
	for path, count := range globalMetrics.pathCounts {
		fmt.Fprintf(w, "http_requests_by_path_total{path=\"%s\"} %d\n", path, count)
	}

	fmt.Fprintf(w, "# HELP http_requests_by_status_total Total requests broken down by HTTP status code\n")
	fmt.Fprintf(w, "# TYPE http_requests_by_status_total counter\n")
	for status, count := range globalMetrics.statusCounts {
		fmt.Fprintf(w, "http_requests_by_status_total{status=\"%d\"} %d\n", status, count)
	}
}

func getEnv(key, defaultVal string) string {
	if val, ok := os.LookupEnv(key); ok && val != "" {
		return val
	}
	return defaultVal
}

func NewRouter() http.Handler {
	mux := http.NewServeMux()
	mux.HandleFunc("/healthz", healthzHandler)
	mux.HandleFunc("/ready", readinessHandler)
	mux.HandleFunc("/api/v1/info", infoHandler)
	mux.HandleFunc("/api/v1/data", dataHandler)
	mux.HandleFunc("/metrics", metricsHandler)

	// Serve Embedded AegisCloud Web Dashboard on "/" and 404 on unknown paths
	mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
		if r.URL.Path != "/" {
			http.NotFound(w, r)
			return
		}
		if strings.Contains(r.Header.Get("Accept"), "application/json") && r.URL.Query().Get("format") == "json" {
			w.Header().Set("Content-Type", "application/json")
			_ = json.NewEncoder(w).Encode(map[string]string{
				"message": "Welcome to AegisCloud Enterprise Multi-Cloud DevSecOps Platform",
				"status":  "ONLINE",
				"docs":    "/api/v1/info",
			})
			return
		}
		w.Header().Set("Content-Type", "text/html; charset=utf-8")
		_, _ = w.Write(dashboardHTML)
	})

	return loggingAndMetricsMiddleware(mux)
}

func main() {
	portStr := getEnv("PORT", "8080")
	port, err := strconv.Atoi(portStr)
	if err != nil {
		port = 8080
	}

	server := &http.Server{
		Addr:         fmt.Sprintf(":%d", port),
		Handler:      NewRouter(),
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt, syscall.SIGTERM)

	go func() {
		log.Printf("Starting core-service on :%d (Cloud: %s, Cluster: %s, Env: %s)", port, cloudProvider, clusterName, environment)
		if err := server.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			log.Fatalf("Server error: %v", err)
		}
	}()

	<-stop
	log.Println("Shutting down core-service gracefully...")

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()

	if err := server.Shutdown(ctx); err != nil {
		log.Fatalf("Server forced to shutdown: %v", err)
	}

	log.Println("Server exiting gracefully.")
}
