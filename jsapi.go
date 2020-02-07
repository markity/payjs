package payjs

import "errors"

// JSAPIInfo JSAPI支付信息表
type JSAPIInfo struct {
	TotalFee   int
	OutTradeNo string
	Body       string
	Attach     string
	NotifyUrl  string
	OpenID     string
}

// check 检查必填项是否为空
func (info *JSAPIInfo) checkEmpty() error {
	if info.TotalFee <= 0 {
		return errors.New("TotalFee must be greater than 0")
	}
	if info.OutTradeNo == "" {
		return errors.New("OutTradeNo cannot be empty")
	}
	if info.OpenID == "" {
		return errors.New("OpenID cannot be empty")
	}
	return nil
}
