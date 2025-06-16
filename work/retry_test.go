package work

import (
	"context"
	"fmt"
	"github.com/vincent78/butil/utils/timeUtil"
	"testing"
	"time"
)

func TestDoneWithRetry(t *testing.T) {
	var txID string
	opName := fmt.Sprintf("do Action %s", timeUtil.NowStr())
	TotalRetryTime := 10 * time.Second
	ctx, cancel := context.WithTimeout(context.Background(), TotalRetryTime)
	defer cancel()

	err := DoneWithRetry(ctx, opName, func() error {
		var attemptErr error
		txID, attemptErr = timeUtil.NowStr(), nil
		return attemptErr
	})

	if err != nil {
		t.Errorf("%s err: %v", opName, err)
	}
	t.Log(txID)
}
