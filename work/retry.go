package work

import (
	"context"
	"fmt"
	"time"

	"github.com/vincent78/butil/logger/logger2/logger"
)

const MaxRetries = 3

// withRetry 通用重试逻辑
func DoneWithRetry(ctx context.Context, operationName string, fn func() error) error {
	var lastErr error
	baseDelay := 500 * time.Millisecond

	for i := 0; i < MaxRetries; i++ {
		select {
		case <-ctx.Done():
			errMsg := fmt.Sprintf("Operation '%s' cancelled via context (attempt %d/%d). Last error if any: %v. Context error: %v", operationName, i+1, MaxRetries, lastErr, ctx.Err())
			logger.Warn(errMsg)
			if lastErr != nil {
				return fmt.Errorf("%s: %w (context error: %v)", errMsg, lastErr, ctx.Err())
			}
			return ctx.Err()
		default:
		}

		logger.Debugf("Attempt %d/%d for operation '%s'...", i+1, MaxRetries, operationName)
		err := fn()
		if err == nil {
			logger.Debugf("Operation '%s' succeeded on attempt %d/%d.", operationName, i+1, MaxRetries)
			return nil
		}
		lastErr = err
		logger.Warnf("Attempt %d/%d for operation '%s' failed: %v", i+1, MaxRetries, operationName, err)

		if i == MaxRetries-1 {
			break
		}

		delay := baseDelay * time.Duration(1<<(i))
		if deadline, ok := ctx.Deadline(); ok {
			if remaining := time.Until(deadline); remaining > 0 && delay > remaining {
				delay = remaining
			}
		}

		if delay <= 0 {
			logger.Warnf("Delay time for retry is <= 0 for operation '%s'. Breaking retry loop.", operationName)
			break
		}

		logger.Debugf("Waiting %v before next attempt for '%s'.", delay, operationName)
		select {
		case <-time.After(delay):
		case <-ctx.Done():
		}
	}

	errMsg := fmt.Sprintf("Operation '%s' failed after %d attempts. Last error: %v", operationName, MaxRetries, lastErr)
	logger.Error(errMsg)
	return fmt.Errorf("operation '%s' failed after %d retries: %w", operationName, MaxRetries, lastErr)
}
