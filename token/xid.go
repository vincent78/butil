package token



import (
	"github.com/rs/xid"
)

type gen string

type ID struct {
	S       string
	Machine string
	PID     uint16
	Counter int32
}

var X gen

func (g *gen) NewID() (string, *ID) {
	guid := xid.New()
	id := &ID{
		S:       guid.String(),
		Machine: string(guid.Machine()),
		PID:     guid.Pid(),
		Counter: guid.Counter(),
	}
	return guid.String(), id
}