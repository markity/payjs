package payjs

import (
	"encoding/json"
	"errors"
)

// ReverseOrderInfo 订单撤销信息表
type ReverseOrderInfo struct {
	PayjsOrderID string // Y PAYJS 平台订单号
}

// checkEmpty 检验必填字段是否为空
func (reverseOrderInfo *ReverseOrderInfo) checkEmpty() error {
	if reverseOrderInfo.PayjsOrderID == "" {
		return errors.New("PayjsOrderID cannot be empty")
	}
	return nil
}

// reverseOrderRequest 订单撤销请求的json结构体
type reverseOrderRequest struct {
	PayjsOrderID string `json:"payjs_order_id"`
	Sign         string `json:"sign"`
}

func (reverseOrderReq *reverseOrderRequest) setSign(mchKey string) {
	reverseOrderReq.Sign = toolSignReq(reverseOrderReq, mchKey)
}

func (reverseOrderReq *reverseOrderRequest) marshal() []byte {
	b, _ := json.Marshal(reverseOrderReq)
	return b
}

// ReverseOrderResponse 订单撤销响应的json结构体
type ReverseOrderResponse struct {
	ReturnCode int    `json:"return_code"` // Y 1:请求成功 0:请求失败, 若失败 mch.Close将返回错误
	ReturnMsg  string `json:"return_msg"`  // Y 返回消息
	Sign       string `json:"sign"`        // Y 数据签名, mch.Close方法会检验签名, 若签名有误则返回错误
}
