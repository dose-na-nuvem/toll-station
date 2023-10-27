package telemetry

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"github.com/stretchr/testify/require"
	"go.opentelemetry.io/otel/attribute"
	semconv "go.opentelemetry.io/otel/semconv/v1.21.0"
)

func Test_newResourceWithBasicAttrs(t *testing.T) {

	// arrange

	// act
	res, err := newResource()
	// assert
	require.NoError(t, err, "o recurso não foi construído")
	require.NotNil(t, res)

	// Verify that the resource has the expected attributes.
	expectedAttributes := make([]attribute.KeyValue, 0)
	expectedAttributes = append(expectedAttributes, semconv.ServiceName(serviceName))
	expectedAttributes = append(expectedAttributes, semconv.ServiceVersion(serviceVersion))

	assert.Subset(t, res.Attributes(), expectedAttributes, "o recurso não tem os attributos básicos definidos")

}
