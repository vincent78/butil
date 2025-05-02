package gcnet

import (
	"fmt"
	"gitee.com/vincent78/gcutil/net/gcnet/proto.pb"
	"gitee.com/vincent78/gcutil/utils/strUtil"
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
