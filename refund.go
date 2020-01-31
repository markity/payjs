package payjs

import (
	"encoding/json"
	"errors"
)

// RefundInfo 退款信息表
type RefundInfo struct {
	PayjsOrderID string // Y PAYJS 平台订单号
}

// checkEmpty 检查必填项是否为空
func (refundInfo *RefundInfo) checkEmpty() error {
	if refundInfo.PayjsOrderID == "" {
		return errors.New("PayjsOrderID cannot be empty")
	}
	return nil
}

type refundRequest struct {
	PayjsOrderID string `json:"payjs_order_id"`
	Sign         string `json:"sign"`
}

func (refundReq *refundRequest) setSign(mchKey string) {
	refundReq.Sign = toolSignReq(refundReq, mchKey)
}

func (refundReq *refundRequest) marshal() []byte {
	b, _ := json.Marshal(refundReq)
	return b
}

type RefundResponse struct {
	ReturnCode int    `json:"return_code"` // Y 1:请求成功 0:请求失败, 若失败则mch.Refund方法返回错误
	ReturnMsg  string `json:"return_msg"`  // Y 返回消息
	Sign       string `json:"sign"`        // Y 数据签名, 若签名错误则mch.Refund方法返回错误
}
