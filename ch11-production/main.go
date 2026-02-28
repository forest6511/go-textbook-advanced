package main

import (
	"context"
	"fmt"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"runtime"
	"syscall"
	"time"
)

var logger = slog.New(slog.NewJSONHandler(os.Stdout, nil))

// healthHandler returns liveness status
func healthHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	fmt.Fprint(w, `{"status":"ok"}`)
}

// Container-aware GOMAXPROCS info (Go 1.25)
func printGOMAXPROCS() {
	procs := runtime.GOMAXPROCS(0)
	numCPU := runtime.NumCPU()
	logger.Info("runtime config",
		slog.Int("GOMAXPROCS", procs),
		slog.Int("NumCPU", numCPU),
		slog.String("note", "Go 1.25: automatically respects cgroup CPU limits"),
	)
}

func main() {
	// Go 1.25: GOMAXPROCS automatically respects cgroup CPU bandwidth limits
	printGOMAXPROCS()

	mux := http.NewServeMux()
	mux.HandleFunc("GET /healthz", healthHandler)
	mux.HandleFunc("GET /readyz", healthHandler)

	srv := &http.Server{
		Addr:    ":8080",
		Handler: mux,
	}

	// Graceful shutdown with signal.NotifyContext
	ctx, stop := signal.NotifyContext(context.Background(), syscall.SIGTERM, syscall.SIGINT)
	defer stop()

	go func() {
		logger.Info("server starting", slog.String("addr", srv.Addr))
		if err := srv.ListenAndServe(); err != nil && err != http.ErrServerClosed {
			logger.Error("server error", slog.Any("err", err))
			os.Exit(1)
		}
	}()

	<-ctx.Done()
	logger.Info("shutting down...")

	shutdownCtx, cancel := context.WithTimeout(context.Background(), 25*time.Second)
	defer cancel()

	if err := srv.Shutdown(shutdownCtx); err != nil {
		logger.Error("shutdown error", slog.Any("err", err))
	}
	logger.Info("shutdown complete")
}
