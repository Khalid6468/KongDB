package utils_test

import (
	"testing"

	"github.com/khalid64/kongdb/pkg/utils"
	"github.com/stretchr/testify/assert"
	"go.uber.org/zap"
)

func TestNewLogger(t *testing.T) {
	tests := []struct {
		name     string
		debug    bool
		expected bool
	}{
		{
			name:     "production logger",
			debug:    false,
			expected: false,
		},
		{
			name:     "debug logger",
			debug:    true,
			expected: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			logger := utils.NewLogger(tt.debug)
			assert.NotNil(t, logger)
			assert.NotNil(t, logger.Logger)
		})
	}
}

func TestLogger_With(t *testing.T) {
	logger := utils.NewLogger(false)

	child := logger.With(zap.String("key", "value"))
	assert.NotNil(t, child)
	assert.NotEqual(t, logger, child)
}

func TestLogger_WithContext(t *testing.T) {
	logger := utils.NewLogger(false)

	ctx := map[string]interface{}{
		"key1": "value1",
		"key2": 123,
	}

	child := logger.WithContext(ctx)
	assert.NotNil(t, child)
	assert.NotEqual(t, logger, child)
}

func TestLogger_Info(t *testing.T) {
	logger := utils.NewLogger(false)

	// This should not panic
	logger.Info("test message")
	logger.Info("test message with fields", zap.String("key", "value"))
}

func TestLogger_Error(t *testing.T) {
	logger := utils.NewLogger(false)

	// This should not panic
	logger.Error("test error")
	logger.Error("test error with fields", zap.String("key", "value"))
}
