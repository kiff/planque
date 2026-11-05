// Command browserworker is the browser worker (kiff/planque#86): the
// headless browser planque drives for a browser grant, run next to planque
// (hosted) or beside it (self-hosted), reached only by planque with the
// shared key.
//
// Environment:
//
//	PLANQUE_WORKER_LISTEN      address to listen on (default :8787)
//	PLANQUE_WORKER_KEY_SHA256  SHA-256, hex, of the key planque presents
//	PLANQUE_CHROME             the browser's path (default: looked up)
//	PLANQUE_WORKER_IDLE        how long an unused session lives (default 2m)
package main

import (
	"context"
	"log/slog"
	"net/http"
	"os"
	"os/signal"
	"syscall"
	"time"

	"github.com/kiff/planque/internal/browserworker"
)

func main() {
	logger := slog.New(slog.NewJSONHandler(os.Stdout, nil))
	listen := os.Getenv("PLANQUE_WORKER_LISTEN")
	if listen == "" {
		listen = ":8787"
	}
	keySum := os.Getenv("PLANQUE_WORKER_KEY_SHA256")
	if len(keySum) != 64 {
		logger.Error("PLANQUE_WORKER_KEY_SHA256 is required: the SHA-256, hex, of planque's key")
		os.Exit(2)
	}
	idle := 2 * time.Minute
	if v := os.Getenv("PLANQUE_WORKER_IDLE"); v != "" {
		d, err := time.ParseDuration(v)
		if err != nil {
			logger.Error("PLANQUE_WORKER_IDLE isn't a duration", slog.String("value", v))
			os.Exit(2)
		}
		idle = d
	}
	srv := &browserworker.Server{Local: &browserworker.Local{Options: browserworker.Options{Chrome: os.Getenv("PLANQUE_CHROME"), Logger: logger}},
		KeySHA256: keySum, Idle: idle, Logger: logger}
	hs := &http.Server{Addr: listen, Handler: srv.Handler(), ReadHeaderTimeout: 10 * time.Second}
	go func() {
		stop := make(chan os.Signal, 1)
		signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)
		<-stop
		ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
		defer cancel()
		srv.Shutdown(ctx)
		_ = hs.Shutdown(ctx)
	}()
	logger.Info("browser worker listening", slog.String("addr", listen))
	if err := hs.ListenAndServe(); err != nil && err != http.ErrServerClosed {
		logger.Error("serve", slog.String("err", err.Error()))
		os.Exit(1)
	}
}
