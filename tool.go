package payjs

import (
	"crypto/md5"
	"encoding/json"
	"fmt"
	"sort"
	"strings"
)

// toolSignReq 为json结构体生成签名
func toolSignReq(structData interface{}, mchKey string) string {
	structDataBytes, _ := json.Marshal(structData)
	kvMap := make(map[string]interface{})
	json.Unmarshal(structDataBytes, &kvMap)
	delete(kvMap, "sign")

	keys := make([]string, 0)
	for k, _ := range kvMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	kvList := make([]string, 0, len(keys)+1)
	for _, k := range keys {
		v := fmt.Sprintf("%v", kvMap[k])
		kvList = append(kvList, fmt.Sprintf("%v=%v", k, v))
	}
	kvList = append(kvList, fmt.Sprintf("key=%v", mchKey))

	sign := strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(kvList, "&")))))
	return sign
}

// toolCheckSignResp 校验响应签名
func toolCheckSignResp(respData []byte, mchKey string) bool {
	kvMap := make(map[string]interface{})
	json.Unmarshal(respData, &kvMap)
	sign, exist := kvMap["sign"]
	// 对于某些特殊情况, 如重复退款, 即使return_code为1也无sign
	if !exist {
		return true
	}
	delete(kvMap, "sign")

	keys := make([]string, 0)
	for k, _ := range kvMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	kvList := make([]string, 0, len(keys)+1)
	for _, k := range keys {
		v := fmt.Sprintf("%v", kvMap[k])
		kvList = append(kvList, fmt.Sprintf("%v=%v", k, v))
	}
	kvList = append(kvList, fmt.Sprintf("key=%v", mchKey))

	realSign := strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(kvList, "&")))))
	if sign != realSign {
		return false
	}
	return true
}
