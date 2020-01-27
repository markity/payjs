package payjs

import (
	"encoding/json"
	"errors"
)

// CloseOrderInfo 订单关闭信息表
type CloseOrderInfo struct {
	PayjsOrderID string // Y PAYJS 平台订单号
}

// check 检验必填字段是否为空
func (closeOrderInfo *CloseOrderInfo) checkEmpty() error {
	if closeOrderInfo.PayjsOrderID == "" {
		return errors.New("PayjsOrderID cannot be empty")
	}
	return nil
}

// closeOrderRequest 订单关闭请求的json结构体
type closeOrderRequest struct {
	PayjsOrderID string `json:"payjs_order_id"`
	Sign         string `json:"sign"`
}

func (closeOrderReq *closeOrderRequest) setSign(mchKey string) {
	closeOrderReq.Sign = toolSignReq(closeOrderReq, mchKey)
}

func (closeOrderReq *closeOrderRequest) marshal() []byte {
	b, _ := json.Marshal(closeOrderReq)
	return b
}

// CloseOrderResponse 订单关闭响应的json结构体
type CloseOrderResponse struct {
	ReturnCode   int    `json:"return_code"`    // Y 1:请求成功 0:请求失败, 若失败 mch.Close将返回错误
	ReturnMsg    string `json:"return_msg"`     // Y 返回消息
	PayjsOrderID string `json:"payjs_order_id"` // Y PAYJS 平台订单号
	Sign         string `json:"sign"`           // Y 数据签名, mch.Close方法会检验签名, 若签名有误则返回错误
}
