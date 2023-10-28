package telemetry

import (
	"context"
	"fmt"
	"net"
	"net/http"
	"testing"

	"github.com/stretchr/testify/require"
	"github.com/tj/assert"
	"go.opentelemetry.io/otel/sdk/resource"
)

func TestNewExporterWithoutError(t *testing.T) {
	// arrange

	// act
	exp, err := newMetricsExporter()
	// assert
	require.NoError(t, err)
	require.NotNil(t, exp)
}

func TestSetupWithoutError(t *testing.T) {
	// arrange
	ctx := context.Background()

	// act
	err := setupMetrics(ctx, resource.Empty())

	// assert
	assert.NoError(t, err)
}

func TestServeMetrics(t *testing.T) {
	// arrange
	// Criate um handle de test handler que incrementa o contador toda vez que é chamado.
	counter := 0
	testHandler := http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		counter++
	})

	port := getFreePortWithFallback(45000)
	endpoint := fmt.Sprintf(":%d", port)

	// Chama  serveMetrics() com o endpoint e o handler.
	go func() {
		err := serveMetrics(endpoint, testHandler)
		if err != nil {
			t.Errorf("serveMetrics() failed: %v", err)
		}
	}()

	// act
	// Make a request to the metrics endpoint.
	resp, err := http.Get("http://localhost" + endpoint + "/metrics")
	if err != nil {
		t.Errorf("GET /metrics failed: %v", err)
	}

	// asserts
	assert.Equal(t, 200, resp.StatusCode, "o endpoint de métricas não está funcionando, verifique o handler")

	// Verify that the counter was incremented.
	if counter != 1 {
		t.Errorf("counter was not incremented: expected 1, got %d", counter)
	}

}

func getFreePortWithFallback(defaultPort int) (port int) {
	var a *net.TCPAddr
	var err error
	if a, _ = net.ResolveTCPAddr("tcp", "localhost:0"); err == nil {
		var l *net.TCPListener
		if l, err = net.ListenTCP("tcp", a); err == nil {
			defer l.Close()
			return l.Addr().(*net.TCPAddr).Port
		}
	}
	return defaultPort
}
