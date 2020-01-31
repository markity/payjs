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
	// 因为调用此方法前, 已经将respData Unmarshal为resp结构体, 证明json数据没有错误, 忽略error
	json.Unmarshal(respData, &kvMap)

	// 获取json数据里的签名
	originSign, ok := kvMap["sign"]
	// 对于某些特殊情况, 如重复退款, 即使return_code为1也无sign
	// 所以如果respData中无sign字段, 则返回true
	if !ok {
		return true
	}

	// 获取keys并排序
	delete(kvMap, "sign")
	keys := make([]string, 0)
	for k, _ := range kvMap {
		keys = append(keys, k)
	}
	sort.Strings(keys)

	// 将key=val追加入kvList中
	kvList := make([]string, 0)
	var walkMap func(string, interface{}, interface{}) = nil
	walkMap = func(perfix string, obj interface{}, ks interface{}) {
		// obj为[]interface{}或者map[string]interface{}类型, ks为[]string或者[]int类型
		switch keys := ks.(type) {
		case []string:
			// 当ks为[]string类型时, obj为map[string]interface{}类型
			object := obj.(map[string]interface{})
			for _, key := range keys {
				v := object[key]
				// 根据签名算法, 值为null的不计入
				if v == nil {
					continue
				}
				// 若v为复合类型, 进入继续遍历
				switch value := v.(type) {
				case []interface{}:
					_perfix := ""
					if perfix == "" {
						_perfix = key
					} else {
						_perfix = fmt.Sprintf("%v[%v]", perfix, key)
					}
					_keys := make([]int, 0, len(value))
					for i := 0; i < len(value); i++ {
						_keys = append(_keys, i)
					}
					walkMap(_perfix, value, _keys)
					continue
				case map[string]interface{}:
					_perfix := ""
					if perfix == "" {
						_perfix = key
					} else {
						_perfix = fmt.Sprintf("%v[%v]", perfix, key)
					}
					_keys := make([]string, 0)
					for k, _ := range value {
						_keys = append(_keys, k)
					}
					sort.Strings(_keys)
					walkMap(_perfix, value, _keys)
					continue
				}
				// 为基础类型
				if perfix == "" {
					kvList = append(kvList, fmt.Sprintf("%v=%v", key, v))
				} else {
					kvList = append(kvList, fmt.Sprintf("%v[%v]=%v", perfix, key, v))
				}
			}
		case []int:
			// keys为[]int类型, perfix不可能为""
			// 当ks为[]int, obj为[]interface{}
			object := obj.([]interface{})
			for _, key := range keys {
				v := object[key]
				if v == nil {
					continue
				}
				switch value := v.(type) {
				// 若v为复合类型, 进入继续遍历
				case []interface{}:
					_perfix := fmt.Sprintf("%v[%v]", perfix, key)
					_keys := make([]int, 0, len(value))
					for i := 0; i < len(value); i++ {
						_keys = append(_keys, i)
					}
					walkMap(_perfix, value, _keys)
					continue
				case map[string]interface{}:
					_perfix := fmt.Sprintf("%v[%v]", perfix, key)
					_keys := make([]string, 0)
					for k, _ := range value {
						_keys = append(_keys, k)
					}
					sort.Strings(_keys)
					walkMap(_perfix, value, _keys)
					continue
				}
				// 为基础类型
				kvList = append(kvList, fmt.Sprintf("%v[%v]=%v", perfix, key, v))
			}
		}
	}
	walkMap("", kvMap, keys)
	kvList = append(kvList, fmt.Sprintf("key=%v", mchKey))

	// 比较签名
	realSign := strings.ToUpper(fmt.Sprintf("%x", md5.Sum([]byte(strings.Join(kvList, "&")))))
	if originSign != realSign {
		return false
	}
	return true
}
