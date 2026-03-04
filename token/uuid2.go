package token

import (
	"fmt"
	"hash/fnv"
	"math/rand"
	"os"
	"runtime"
	"strconv"
	"time"

	"github.com/kjk/betterguid"
	"github.com/oklog/ulid"
	"github.com/rs/xid"
	"github.com/segmentio/ksuid"
	"github.com/sony/sonyflake/v2"
	"github.com/spf13/cast"
	"github.com/vincent78/butil/net/common"
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
	flake := NewSonyFlake(sonyflake.Settings{})
	id, err := flake.GenInt64()
	if err != nil {
		panic(err)
	}
	return strconv.FormatUint(cast.ToUint64(id), 16)
}

const (
	nodeBits  = 10 // 节点 ID：最多 1024 个实例
	maxNodeID = (1 << nodeBits) - 1
)

var machineID int

type SonyFlake struct {
	sf *sonyflake.Sonyflake
}

func GenerateUserID() (uint64, error) {
	idGenerator := NewSonyFlake()
	id, err := idGenerator.GenInt64()
	if err != nil {
		return 0, err
	}
	return uint64(id), nil
}

func NewSonyFlake(settings ...sonyflake.Settings) *SonyFlake {
	machineID = int(genNodeID())
	setting := sonyflake.Settings{
		MachineID: func() (int, error) {
			return machineID, nil
		},
	}
	if len(settings) > 0 {
		setting = settings[0]
	}

	sf, err := sonyflake.New(setting)
	if err != nil {
		return nil
	}
	return &SonyFlake{sf}
}

func (s *SonyFlake) GenInt64() (int64, error) {
	return s.sf.NextID()
}

// genNodeID 自动生成容器内唯一的 NodeID
func genNodeID() int64 {
	ip := common.LocalIP()
	pid := os.Getpid()
	cpu := runtime.NumCPU()
	now := time.Now().UnixNano()
	hostname := os.Getenv("HOSTNAME")

	// 多实例不同：IP + PID + 纳秒启动时间 + CPU核数
	raw := fmt.Sprintf("%s-%d-%d-%d-%s", ip, pid, cpu, now, hostname)

	h := fnv.New32a()
	h.Write([]byte(raw))

	return int64(h.Sum32()) & maxNodeID
}
