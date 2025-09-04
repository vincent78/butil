package gcgin

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"strconv"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

// ParseJSON 解析请求JSON
func ParseJSON(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindJSON(obj); err != nil {
		return errors.New(fmt.Sprintf("解析请求参数发生错误 - %s", err.Error()))
	}
	return nil
}

// ParseQuery 解析Query参数
func ParseQuery(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindQuery(obj); err != nil {
		return errors.New(fmt.Sprintf("解析请求参数发生错误 - %s", err.Error()))
	}
	return nil
}

// ParseForm 解析Form请求
func ParseForm(c *gin.Context, obj interface{}) error {
	if err := c.ShouldBindWith(obj, binding.Form); err != nil {
		return errors.New(fmt.Sprintf("解析请求参数发生错误 - %s", err.Error()))
	}
	return nil
}

func ParseParam(c *gin.Context, param string) string {
	return c.Param(param)
}

func ParseParamUint64(c *gin.Context, param string) uint64 {
	v := ParseParam(c, param)
	vID, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		// log.Printf("解析路由路径参数[%s]错误: %s", param, err.Error())
	}
	return vID
}

func ParseParamUint(c *gin.Context, param string) uint {
	v := ParseParam(c, param)
	vID, err := strconv.ParseUint(v, 10, 64)
	if err != nil {
		// log.Printf("解析路由路径参数[%s]错误: %s", param, err.Error())
	}
	return uint(vID)
}

func ParseParamInt64(c *gin.Context, param string) int64 {
	v := ParseParam(c, param)
	vID, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		// log.Printf("解析路由路径参数[%s]错误: %s", param, err.Error())
	}
	return vID
}

func ParseParamInt(c *gin.Context, param string) int {
	v := ParseParam(c, param)
	vID, err := strconv.ParseInt(v, 10, 64)
	if err != nil {
		// log.Printf("解析路由路径参数[%s]错误: %s", param, err.Error())
	}
	return int(vID)
}

func GetRawDataToMapUseNumber(c *gin.Context) (map[string]interface{}, error) {
	req := c.Request
	d := json.NewDecoder(req.Body)
	d.UseNumber()
	var x interface{}
	if err := d.Decode(&x); err != nil {
		// log.Errorf(err, "decode json error: %s", err.Error())
	}
	maps := x.(map[string]interface{})
	removeTraceCond(maps)
	return maps, nil
}

func GetRawDataToMap(c *gin.Context) (map[string]interface{}, error) {
	var condition = make(map[string]interface{})

	b, err := c.GetRawData()
	if err != nil {
		return nil, err
	}

	if b == nil || len(b) <= 0 {
		return condition, nil
	}
	err = json.Unmarshal(b, &condition)
	if err != nil {
		return nil, err
	}

	removeTraceCond(condition)

	return condition, nil
}

// removeTraceCond 移除客户端的跟踪条件
// TODO: 需重构在拦截器处理跟踪条件(appBuild, deviceName等)
func removeTraceCond(cond map[string]interface{}) {
	keys := []string{"appBuild", "deviceName"}

	for _, key := range keys {
		if _, ok := cond[key]; ok {
			delete(cond, key)
		}
	}

}

func GetRawDataToString(c *gin.Context) (string, error) {
	req := c.Request
	result, err := ioutil.ReadAll(req.Body)
	if err != nil {
		return "", err
	} else {
		return bytes.NewBuffer(result).String(), nil
	}
}
