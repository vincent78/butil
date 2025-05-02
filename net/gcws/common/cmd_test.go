package common

import (
	"gitee.com/vincent78/gcutil/utils/timeUtil"
	"testing"
	"time"
)

func TestPingCmdCreate(t *testing.T) {
	cmd := NewWSCmdByCategory(CmdCategorySys, CmdNamePing)
	cmd.Data = timeUtil.NowStr()
	t.Logf("ping cmd : %v", cmd.String())
}

func TestLoginCmdCreate(t *testing.T) {
	//cmd := NewWSCmdByCategory(CmdCategorySys, CmdNameLogin)
	//cmd.Data = map[string]string{
	//	""
	//}
}

func TestCmdRespSuccess(t *testing.T) {
	cmd := NewWSCmdByCategory(CmdCategorySys, CmdNamePing)
	cmd.Data = timeUtil.NowStr()
	t.Logf("cmd : %v", cmd.String())
	time.Sleep(time.Second)
	cmdResp := cmd.Success(timeUtil.NowStr())
	t.Logf("resp : %v", cmdResp.String())
}
