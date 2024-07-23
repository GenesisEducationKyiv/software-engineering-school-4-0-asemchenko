package util

import (
	"go.uber.org/zap"
	"io"
)

func CloseBodyWithErrorHandling(Body io.ReadCloser) {
	err := Body.Close()
	if err != nil {
		zap.S().Warnf("Failed to close HTTP request body %v", err)
	}
}

func WarnOnError(err error) {
	if err != nil {
		zap.S().Warn(err)
	}
}
