package payjs

import (
	"encoding/json"
	"errors"
)

// MicroPayInfo 付款码支付信息表
type MicroPayInfo struct {
	TotalFee   int    // Y 金额。单位：分
	OutTradeNo string // Y 用户端自主生成的订单号

	Body   string // N 订单标题
	Attach string // N 用户自定义数据，在notify的时候会原样返回

	AuthCode string // Y扫码支付授权码，设备读取用户微信中的条码或者二维码信息(注：用户刷卡条形码规则：18位纯数字，以10、11、12、13、14、15开头)
}

// checkEmpty 检查必填项是否为空
func (microPayInfo *MicroPayInfo) checkEmpty() error {
	if microPayInfo.TotalFee <= 0 {
		return errors.New("TotalFee must be greater than 0")
	}
	if microPayInfo.OutTradeNo == "" {
		return errors.New("OutTradeNo cannot be empty")
	}
	if microPayInfo.AuthCode == "" {
		return errors.New("AuthCode cannot be empty")
	}
	return nil
}

// microPayRequest 付款码支付请求的json结构体
type microPayRequest struct {
	MchID      string `json:"mchid"`
	TotalFee   int    `json:"total_fee"`
	OutTradeNo string `json:"out_trade_no"`

	Body   string `json:"body,omitempty"`
	Attach string `json:"attach,omitempty"`

	AuthCode string `json:"auth_code"`
	Sign     string `json:"sign"`
}

func (microPayReq *microPayRequest) setSign(mchKey string) {
	microPayReq.Sign = toolSignReq(microPayReq, mchKey)
}

func (microPayReq *microPayRequest) marshal() []byte {
	b, _ := json.Marshal(microPayReq)
	return b
}

// MicroPayResponse 撤销订单响应
type MicroPayResponse struct {
	ReturnCode   int    `json:"return_code"`    // Y 1:请求成功，0:请求失败
	ReturnMsg    string `json:"return_msg"`     // Y 返回消息
	PayjsOrderID string `json:"payjs_order_id"` // Y PAYJS 平台订单号
	Sign         string `json:"sign"`           // Y 数据签名
}
