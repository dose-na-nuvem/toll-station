package telemetry

import (
	"context"
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
)

func TestNew(t *testing.T) {
	// arrange
	ctx := context.Background()

	// act
	telemetry, err := New(ctx)

	// assert
	require.NoError(t, err)
	require.NotNil(t, telemetry)
	assert.NotNil(t, telemetry.TrafficCounter, "o contador é invalido")
	assert.NotNil(t, telemetry.GateHistogram, "o contador é invalido")

}
