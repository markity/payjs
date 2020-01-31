package payjs

import (
	"encoding/json"
	"errors"
)

// NativeInfo 扫码支付信息表
type NativeInfo struct {
	TotalFee   int    // Y 金额。单位：分
	OutTradeNo string // Y 用户端自主生成的订单号

	Body      string // N 订单标题
	Attach    string // N 用户自定义数据，在notify的时候会原样返回
	NotifyUrl string // N 接收微信支付异步通知的回调地址。必须为可直接访问的URL，不能带参数、session验证、csrf验证。留空则不通知
}

// check 检查必填项是否为空
func (info *NativeInfo) checkEmpty() error {
	if info.TotalFee <= 0 {
		return errors.New("TotalFee must be greater than 0")
	}
	if info.OutTradeNo == "" {
		return errors.New("OutTradeNo cannot be empty")
	}
	return nil
}

// nativeRequest Native请求的json结构体
type nativeRequest struct {
	MchID      string `json:"mchid"`
	TotalFee   int    `json:"total_fee"`
	OutTradeNo string `json:"out_trade_no"`

	Body      string `json:"body,omitempty"`
	Attach    string `json:"attach,omitempty"`
	NotifyUrl string `json:"notify_url,omitempty"`

	Sign string `json:"sign"`
}

// setSign 设置签名
func (nativeReq *nativeRequest) setSign(mchKey string) {
	nativeReq.Sign = toolSignReq(nativeReq, mchKey)
}

// marshal 结构体编码为json
func (nativeReq *nativeRequest) marshal() []byte {
	b, _ := json.Marshal(nativeReq)
	return b
}

// NativeResponse 扫码支付请求返回值
type NativeResponse struct {
	ReturnCode   int    `json:"return_code"`    // Y 1:请求成功,0:请求失败.若请求失败, mch.Native方法将返回错误
	ReturnMsg    string `json:"return_msg"`     // Y 返回消息
	PayjsOrderID string `json:"payjs_order_id"` // Y PAYJS 平台订单号
	Qrcode       string `json:"qrcode"`         // Y 二维码图片地址
	CodeUrl      string `json:"code_url"`       // Y 可将该参数生成二维码展示出来进行扫码支付(有效期2小时)
	Sign         string `json:"sign"`           // Y 数据签名, 用于验证请求的合法性, 和校验请求信息正误.若签名错误, mch.Native方法将返回错误
}
