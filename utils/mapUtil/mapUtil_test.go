package mapUtil

import (
	"encoding/json"
	"testing"

	"github.com/vincent78/butil/logger/logger1"
)

func TestGetObjByAnchorFromMap(b *testing.T) {
	level := []string{"test1", "test2"}
	str := `{"test1":{"test2":"this is the test2"}}`
	var obj map[string]any
	_ = json.Unmarshal([]byte(str), &obj)
	logger1.Debug("--- %v", ObjByAnchorFromMap(obj, level))
}
