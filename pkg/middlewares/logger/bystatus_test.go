package logger

import (
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestLogFilteredStatusCode_WhenStatusIs5xx_ShouldLogAsError(t *testing.T) {
	// Given
	statusCode := 503
	prefix := "test"

	// When
	logFilteredStatusCode(statusCode, prefix)

	// Then
	assert.GreaterOrEqual(t, statusCode, 500)
	assert.Less(t, statusCode, 600)
}

func TestLogFilteredStatusCode_WhenStatusIs4xx_ShouldLogAsWarning(t *testing.T) {
	// Given
	statusCode := 404
	prefix := "test"

	// When
	logFilteredStatusCode(statusCode, prefix)

	// Then
	assert.GreaterOrEqual(t, statusCode, 400)
	assert.Less(t, statusCode, 500)
}

func TestLogFilteredStatusCode_WhenStatusIs2xx_ShouldLogAsInfo(t *testing.T) {
	// Given
	statusCode := 200
	prefix := "test"

	// When
	logFilteredStatusCode(statusCode, prefix)

	// Then
	assert.GreaterOrEqual(t, statusCode, 200)
	assert.Less(t, statusCode, 300)
}
