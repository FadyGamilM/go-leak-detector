package goleak

import (
	"context"
	"net/http"
	"time"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

// starts a Prometheus metrics server.
func (m *GoLeakMonitor) StartPrometheusServer(ctx context.Context, addr string) error {
	http.Handle("/metrics", promhttp.HandlerFor(m.monitor.Registry(), promhttp.HandlerOpts{}))
	server := &http.Server{Addr: addr}
	go func() {
		<-ctx.Done()
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		server.Shutdown(ctx)
	}()
	return server.ListenAndServe()
}
