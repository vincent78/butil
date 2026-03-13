package tracing

import (
	"fmt"
	"testing"

	"github.com/vincent78/butil/token"
)

func TestStartSpanWithCustomTraceID(t *testing.T) {
	cp, err := InitJaeger()
	if err != nil {
		t.Fatal(err)
	}
	defer cp()

	//trace id
	traceId := token.NewTraceId()

	ctxSpan, span := StartSpanWithCustomTraceID("service.MatchFilterList", traceId)
	defer span.End()

	endGetSports := StartTimedSpan(ctxSpan, "db.GetSports", map[string]string{
		"lan":      "en",
		"sportIds": fmt.Sprintf("%v", []string{"sr:sport:1", "sr:sport:2", "sr:sport:109", "sr:sport:111"}),
	})
	defer endGetSports()

}
