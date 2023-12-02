/*
@Created by : 2021/1/11 15:20
@Author : hll
@Descripition :
*/
package common

import (
	"bytes"
	"fmt"
	"strconv"
	"strings"
)

// type ResponseServer struct {
// 	Code string      `json:"code"`
// 	Msg  string      `json:"msg"`
// 	Data interface{} `json:"data,omitempty"`
// }

// func NewResponseServer() *ResponseServer {
// 	return &ResponseServer{
// 		Code: "0",
// 		Msg:  "ok",
// 	}
// }

// func (this *ResponseServer) ResDBErr(err error) {
// 	this.Msg = fmt.Sprintf("%s, err:%s", GetErrMsg(RECODE_DB_OPERFAILD), err)
// 	this.Code = RECODE_DB_OPERFAILD
// }

// yaml文件需要替换字段
type YamlReplace struct {
	Old string
	New interface{}
}

// 替换变量，得到最终的yaml文件内容
func ReplaceVar(content []byte, s []YamlReplace) (result []byte) {
	// TODO：为什么需要重新赋值？外部引用此方法又重新转换成了字符串，为啥方法不直接返回字符串？
	result = content
	for _, v := range s {
		switch new := v.New.(type) {
		case string:
			result = bytes.ReplaceAll(result, []byte(v.Old), []byte(new))
		case int:
			result = bytes.ReplaceAll(result, []byte(v.Old), []byte(strconv.Itoa(new)))
		case int64:
			result = bytes.ReplaceAll(result, []byte(v.Old), []byte(strconv.FormatInt(new, 10)))
		case float64: //保留小数点后三位
			result = bytes.ReplaceAll(result, []byte(v.Old), []byte(strconv.FormatFloat(new, 'f', 3, 64)))
		}
	}
	return result
}

// 切割字符串
func SplitStr(str string, num uint) string {
	if len(str) < int(num) {
		return str
	}
	// 010200
	var splited string
	for {
		splited = fmt.Sprintf("%s_%s", splited, str[:num])
		if len(str[num:]) < int(num) {
			return strings.Trim(splited, "_")
		}
		str = str[num:]
	}
}
