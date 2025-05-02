package fileUtil

import (
	"testing"
)

func TestGetCurrentAbPath(t *testing.T) {
	t.Log("getTmpDir（当前系统临时目录） = ", GetTmpDir())
	t.Log("getCurrentAbPathByExecutable（仅支持go build） = ", GetCurrentAbPathByExecutable())
	t.Log("getCurrentAbPathByCaller（仅支持go run） = ", GetCurrentAbPathByCaller())
	t.Log("getCurrentAbPath（最终方案-全兼容） = ", GetCurrentAbPath())

}

func TestGetCurrentProjectPath(t *testing.T) {
	t.Log("GetCurrentProjectPath ", GetCurrentProjectPath())
}

func TestCheckFileExist(t *testing.T) {
	r := CheckFileExist(Join(GetCurrentAbPath(), "pathUtil_test.go"))
	t.Logf("checkFileExist (pathUtil_test.go) : %v", r)
}

func TestGetPrePath(t *testing.T) {
	p := GetCurrentAbPath()
	t.Log("currentPath:", p)
	t.Log("PrePath", GetPrePath(p))

	p = "/"
	t.Log("currentPath:", p)
	t.Log("PrePath", GetPrePath(p))

	p = "/test"
	t.Log("currentPath:", p)
	t.Log("PrePath", GetPrePath(p))

	p = "/test/test1/test2"
	t.Log("currentPath:", p)
	t.Log("PrePath", GetPrePath(p))

	p = "/test/test1/test2/testfile.go"
	t.Log("currentPath:", p)
	t.Log("PrePath", GetPrePath(p))
}
