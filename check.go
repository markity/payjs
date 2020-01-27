package payjs

import (
	"encoding/json"
	"errors"
)

type CheckOrderInfo struct {
	PayjsOrderID string // Y PAYJS 平台订单号
}

func (checkOrderInfo *CheckOrderInfo) checkEmpty() error {
	if checkOrderInfo.PayjsOrderID == "" {
		return errors.New("PayjsOrderId cannot be empty")
	}
	return nil
}

type checkOrderRequest struct {
	PayjsOrderID string `json:"payjs_order_id"`
	Sign         string `json:"sign"`
}

func (checkOrderReq *checkOrderRequest) setSign(mchKey string) {
	checkOrderReq.Sign = toolSignReq(checkOrderReq, mchKey)
}

func (checkOrderReq *checkOrderRequest) marshal() []byte {
	b, _ := json.Marshal(checkOrderReq)
	return b
}

type CheckOrderResponse struct {
	ReturnCode    int    `json:"return_code"`    // Y 1:请求成功 0:请求失败, 若失败mch.CheckOrder方法会返回错误
	MchID         string `json:"mchid"`          // Y PAYJS 平台商户号
	OutTradeNo    string `json:"out_trade_no"`   // Y 用户端订单号
	PayjsOrderID  string `json:"payjs_order_id"` // Y PAYJS 订单号
	TransactionID string `json:"transaction_id"` // N 微信显示订单号
	Status        int    `json:"status"`         // Y 0：未支付，1：支付成功
	Openid        string `json:"openid"`         // N 支付用户的 OPENID(如果支付则有此字段)
	TotalFee      int    `json:"total_fee"`      // N 订单金额
	PaidTime      string `json:"paid_time"`      // N 订单支付时间(如果支付则有此字段)
	Attach        string `json:"attach"`         // N 用户自定义数据
	Sign          string `json:"sign"`           // Y 数据签名 详见签名算法, 若签名有误, mch.CheckOrder方法会返回错误
}
