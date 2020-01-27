package payjs

import (
	"bytes"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
)

// NewMch 新建商户
func NewMch(mchID string, mchKey string) *mch {
	if mchID == "" {
		panic(errors.New("mchID cannot be empty"))
	}
	if mchKey == "" {
		panic(errors.New("mchKey cannot be empty"))
	}
	return &mch{mchID: mchID, mchKey: mchKey}
}

type mch struct {
	mchID  string
	mchKey string
}

// Native 发送扫码支付请求
func (m *mch) Native(nativeInfo NativeInfo) (*NativeResponse, error) {
	if err := nativeInfo.checkEmpty(); err != nil {
		return nil, err
	}
	nativeReq := &nativeRequest{MchID: m.mchID, TotalFee: nativeInfo.TotalFee, OutTradeNo: nativeInfo.OutTradeNo, Body: nativeInfo.Body, Attach: nativeInfo.Attach, NotifyUrl: nativeInfo.NotifyUrl}
	nativeReq.setSign(m.mchKey)
	nativeReqBytes := nativeReq.marshal()
	resp, err := http.DefaultClient.Post("https://payjs.cn/api/native", "application/json", bytes.NewReader(nativeReqBytes))
	if err != nil {
		return nil, err
	}

	nativeResp := &NativeResponse{}
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(respBodyBytes, nativeResp); err != nil {
		return nil, err
	}

	// 判断请求是否成功
	if nativeResp.ReturnCode != 1 {
		return nil, errors.New(fmt.Sprintf("failed: return_code:%v, return_msg:%v", nativeResp.ReturnCode, nativeResp.ReturnMsg))
	}

	// 检验签名
	if !toolCheckSignResp(respBodyBytes, m.mchKey) {
		return nil, errors.New("response signature is wrong")
	}

	return nativeResp, nil
}

// CloseOrder 关闭未支付的订单
func (m *mch) CloseOrder(closeOrderInfo CloseOrderInfo) (*CloseOrderResponse, error) {
	if err := closeOrderInfo.checkEmpty(); err != nil {
		return nil, err
	}
	closeOrderReq := &closeOrderRequest{PayjsOrderID: closeOrderInfo.PayjsOrderID}
	closeOrderReq.setSign(m.mchKey)
	b := closeOrderReq.marshal()
	resp, err := http.DefaultClient.Post("https://payjs.cn/api/close", "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	closeOrderResp := &CloseOrderResponse{}
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(respBodyBytes, closeOrderResp); err != nil {
		return nil, err
	}

	// 判断请求是否成功
	if closeOrderResp.ReturnCode != 1 {
		return nil, errors.New(fmt.Sprintf("failed: return_code:%v, return_msg:%v", closeOrderResp.ReturnCode, closeOrderResp.ReturnMsg))
	}

	// 检验签名
	if !toolCheckSignResp(respBodyBytes, m.mchKey) {
		return nil, errors.New("response signature is wrong")
	}

	return closeOrderResp, nil
}

// Refund 对于已支付的订单退款
func (m *mch) Refund(refundInfo RefundInfo) (*RefundResponse, error) {
	if err := refundInfo.checkEmpty(); err != nil {
		return nil, err
	}
	refundReq := &refundRequest{PayjsOrderID: refundInfo.PayjsOrderID}
	refundReq.setSign(m.mchKey)
	b := refundReq.marshal()
	resp, err := http.DefaultClient.Post("https://payjs.cn/api/refund", "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	refundResp := &RefundResponse{}
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(respBodyBytes, refundResp); err != nil {
		return nil, err
	}

	// 判断请求是否成功
	if refundResp.ReturnCode != 1 {
		return nil, errors.New(fmt.Sprintf("failed: return_code:%v, return_msg:%v", refundResp.ReturnCode, refundResp.ReturnMsg))
	}

	// 检验签名
	if !toolCheckSignResp(respBodyBytes, m.mchKey) {
		return nil, errors.New("response signature is wrong")
	}

	return refundResp, nil
}

// CheckOrder 发起订单查询来确认订单状态
func (m *mch) CheckOrder(checkOrderInfo CheckOrderInfo) (*CheckOrderResponse, error) {
	if err := checkOrderInfo.checkEmpty(); err != nil {
		return nil, err
	}
	checkOrderReq := &checkOrderRequest{PayjsOrderID: checkOrderInfo.PayjsOrderID}
	checkOrderReq.setSign(m.mchKey)
	b := checkOrderReq.marshal()
	resp, err := http.DefaultClient.Post("https://payjs.cn/api/check", "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	checkOrderResp := &CheckOrderResponse{}
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(respBodyBytes, checkOrderResp); err != nil {
		return nil, err
	}

	// 判断请求是否成功
	if checkOrderResp.ReturnCode != 1 {
		return nil, errors.New(fmt.Sprintf("failed: return_code:%v", checkOrderResp.ReturnCode))
	}

	// 检验签名
	if !toolCheckSignResp(respBodyBytes, m.mchKey) {
		return nil, errors.New("response signature is wrong")
	}

	return checkOrderResp, nil
}

// MicroPay 发送付款码支付请求
func (m *mch) MicroPay(microPayInfo MicroPayInfo) (*MicroPayResponse, error) {
	if err := microPayInfo.checkEmpty(); err != nil {
		return nil, err
	}
	microPayReq := &microPayRequest{MchID: m.mchID, TotalFee: microPayInfo.TotalFee, OutTradeNo: microPayInfo.OutTradeNo, Body: microPayInfo.Body, Attach: microPayInfo.Attach, AuthCode: microPayInfo.AuthCode}
	microPayReq.setSign(m.mchKey)
	nativeReqBytes := microPayReq.marshal()
	resp, err := http.DefaultClient.Post("https://payjs.cn/api/micropay", "application/json", bytes.NewReader(nativeReqBytes))
	if err != nil {
		return nil, err
	}

	mircoPayResp := &MicroPayResponse{}
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(respBodyBytes, mircoPayResp); err != nil {
		return nil, err
	}

	// 判断请求是否成功
	if mircoPayResp.ReturnMsg != "需要用户输入支付密码" {
		return nil, errors.New(fmt.Sprintf("failed: return_code:%v, return_msg:%v", mircoPayResp.ReturnCode, mircoPayResp.ReturnMsg))
	}

	// 检验签名
	if !toolCheckSignResp(respBodyBytes, m.mchKey) {
		return nil, errors.New("response signature is wrong")
	}

	return mircoPayResp, nil
}

// ReverseOrder 撤销未支付的订单
func (m *mch) ReverseOrder(reverseOrderInfo ReverseOrderInfo) (*ReverseOrderResponse, error) {
	if err := reverseOrderInfo.checkEmpty(); err != nil {
		return nil, err
	}
	reverseOrderReq := &reverseOrderRequest{PayjsOrderID: reverseOrderInfo.PayjsOrderID}
	reverseOrderReq.setSign(m.mchKey)
	b := reverseOrderReq.marshal()
	resp, err := http.DefaultClient.Post("https://payjs.cn/api/reverse", "application/json", bytes.NewReader(b))
	if err != nil {
		return nil, err
	}

	reverseOrderResp := &ReverseOrderResponse{}
	respBodyBytes, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return nil, err
	}
	if err := json.Unmarshal(respBodyBytes, reverseOrderResp); err != nil {
		return nil, err
	}

	// 校验签名
	if reverseOrderResp.ReturnCode != 1 {
		return nil, errors.New(fmt.Sprintf("failed: return_code:%v, return_msg:%v", reverseOrderResp.ReturnCode, reverseOrderResp.ReturnMsg))
	}

	// 检验签名
	if !toolCheckSignResp(respBodyBytes, m.mchKey) {
		return nil, errors.New("response signature is wrong")
	}

	return reverseOrderResp, nil
}
