package token

import (
	"crypto/sha256"
	"encoding/hex"

	"github.com/google/uuid"
)

func NewTraceId() string {
	traceId := uuid.New().String()
	hash := sha256.Sum256([]byte(traceId))
	actualTraceID := hex.EncodeToString(hash[:16])
	return actualTraceID
}
