package mapUtil

import (
	"encoding/json"
	"gitee.com/vincent78/gcutil/logger"
	"testing"
)

func TestGetObjByAnchorFromMap(b *testing.T) {
	level := []string{"test1", "test2"}
	str := `{"test1":{"test2":"this is the test2"}}`
	var obj map[string]interface{}
	_ = json.Unmarshal([]byte(str), &obj)
	logger.Debug("--- %v", ObjByAnchorFromMap(obj, level))
}
