package gcnet

import (
	"fmt"
	"github.com/vincent78/butil/net/gcnet/proto.pb"
	"github.com/vincent78/butil/utils/strUtil"
	"google.golang.org/protobuf/proto"
	"testing"
)

func TestProto(t *testing.T) {
	cmd := &gcnet.CmdRequest{
		Cmd:     "echo",
		Payload: "hello",
	}
	r, e := proto.Marshal(cmd)
	if e != nil {
		_ = fmt.Errorf("marshal error: %v", e)
	} else {
		fmt.Printf("the marshal :%v\n", r)
	}
	fmt.Println("")
	var tmp gcnet.CmdRequest
	e = proto.Unmarshal(r, &tmp)
	if e != nil {
		_ = fmt.Errorf("unMarshal error: %v", e)
	} else {
		fmt.Printf("the unMarshal :%v\n", strUtil.ToJsonStr(tmp))
	}

}
