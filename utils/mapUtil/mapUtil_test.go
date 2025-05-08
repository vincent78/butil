package mapUtil

import (
	"encoding/json"
	"github.com/vincent78/butil/logger/logger1"
	"testing"
)

func TestGetObjByAnchorFromMap(b *testing.T) {
	level := []string{"test1", "test2"}
	str := `{"test1":{"test2":"this is the test2"}}`
	var obj map[string]interface{}
	_ = json.Unmarshal([]byte(str), &obj)
	logger1.Debug("--- %v", ObjByAnchorFromMap(obj, level))
}
