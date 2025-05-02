package token

import (
	"github.com/chilts/sid"
	"github.com/kjk/betterguid"
	"github.com/oklog/ulid"
	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
	"github.com/sony/sonyflake"
	"github.com/vincent78/butil/logger"
	"math/rand"
	"strconv"
	"time"
)

func GenXid() string {
	return xid.New().String()
}

func GenKsuid() string {
	return ksuid.New().String()
}

func GenBetterGUID() string {
	return betterguid.New()
}

func GenUlid() string {
	t := time.Now().UTC()
	entropy := rand.New(rand.NewSource(t.UnixNano()))
	id := ulid.MustNew(ulid.Timestamp(t), entropy)
	return id.String()
}

func GenSonyflake() string {
	flake := sonyflake.NewSonyflake(sonyflake.Settings{})
	id, err := flake.NextID()
	if err != nil {
		logger.Error("flake.NextID() failed with %s\n", err)
	}
	return strconv.FormatUint(id, 16)
}

func GenSid() string {
	return sid.Id()
}
