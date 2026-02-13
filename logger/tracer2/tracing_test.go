package tracing

import (
	"crypto/sha256"
	"encoding/hex"
	"fmt"
	"testing"

	"github.com/google/uuid"
)

func TestStartSpanWithCustomTraceID(t *testing.T) {
	cp, err := InitJaeger()
	if err != nil {
		t.Fatal(err)
	}
	defer cp()

	//trace id
	traceId := uuid.New().String()
	// TraceID(32位小写十六进制):sha256 前16字节(Jeager)
	hash := sha256.Sum256([]byte(traceId))
	traceId = hex.EncodeToString(hash[:16])

	ctxSpan, span := StartSpanWithCustomTraceID("service.MatchFilterList", traceId)
	defer span.End()

	endGetSports := StartTimedSpan(ctxSpan, "db.GetSports", map[string]string{
		"lan":      "en",
		"sportIds": fmt.Sprintf("%v", []string{"sr:sport:1", "sr:sport:2", "sr:sport:109", "sr:sport:111"}),
	})
	defer endGetSports()

}
