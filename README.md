# PayJS SDK for Go

## 当前实现

- [x] 扫码支付
- [x] 查询订单
- [x] 关闭订单
- [x] 退款
- [ ] 付款码支付
- [ ] 人脸支付
- [ ] 撤销订单
- [ ] 分账支付

## 基本使用

- [扫码支付(native)](#扫码支付(native))
- [查询订单(check)](#查询订单(check))
- [关闭订单(close)](#关闭订单(close))
- [退款(refund)](#退款(refund))

### 扫码支付(native)

```go
package main

import (
	"fmt"
	"payjs/payjs"
)

func main() {
    mch := payjs.NewMch("your-mchid", "your-mchkey")
    // 说明: ReturnCode不为1 或 签名错误均视为error
	nativeResp, err := mch.Native(payjs.NativeInfo{TotalFee: 10, OutTradeNo: "2020_1_27_001", Body: "支付测试"})
	if err != nil {
		fmt.Printf("失败: %v\n", err)
		return
	}

	fmt.Println(nativeResp.ReturnCode)
	fmt.Println(nativeResp.ReturnMsg)
	fmt.Println(nativeResp.PayjsOrderID)
	fmt.Println(nativeResp.OutTradeNo)
	fmt.Println(nativeResp.TotalFee)
	fmt.Println(nativeResp.Qrcode)
	fmt.Println(nativeResp.CodeUrl)
	fmt.Println(nativeResp.Sign)
}
```

### 查询订单(check)

```go
package main

import (
	"fmt"
	"payjs/payjs"
)

func main() {
    mch := payjs.NewMch("your-mchid", "your-mchkey")
    // 说明: ReturnCode不为1 或 签名错误均视为error
	checkResp, err := mch.CheckOrder(payjs.CheckOrderInfo{PayjsOrderID: "order-id"})
	if err != nil {
		fmt.Printf("失败: %v\n", err)
		return
	}

	fmt.Println(checkResp.ReturnCode)
	fmt.Println(checkResp.MchID)
	fmt.Println(checkResp.OutTradeNo)
	fmt.Println(checkResp.PayjsOrderID)
	fmt.Println(checkResp.TransactionID)
	fmt.Println(checkResp.Status)
	fmt.Println(checkResp.Openid)
	fmt.Println(checkResp.TotalFee)
	fmt.Println(checkResp.PaidTime)
	fmt.Println(checkResp.Attach)
	fmt.Println(checkResp.Sign)
}
```

### 关闭订单(close)

```go
package main

import (
	"fmt"
	"payjs/payjs"
)

func main() {
    mch := payjs.NewMch("your-mchid", "your-mchkey")
    // 说明: ReturnCode不为1 或 签名错误均视为error
	closeResp, err := mch.CloseOrder(payjs.CloseOrderInfo{PayjsOrderID: "order-id"})
	if err != nil {
		fmt.Printf("失败: %v\n", err)
		return
	}

	fmt.Println(closeResp.ReturnCode)
	fmt.Println(closeResp.ReturnMsg)
	fmt.Println(closeResp.PayjsOrderID)
	fmt.Println(closeResp.Sign)
}

```

### 退款(refund)

```go
package main

import (
	"fmt"
	"payjs/payjs"
)

func main() {
    mch := payjs.NewMch("your-mchid", "your-mchkey")
    // 说明: ReturnCode不为1 或 签名错误均视为error
	refundResp, err := mch.Refund(payjs.RefundInfo{PayjsOrderID: "order-id"})
	if err != nil {
		fmt.Printf("失败: %v\n", err)
		return
	}

	fmt.Println(refundResp.ReturnCode)
	fmt.Println(refundResp.ReturnMsg)
	fmt.Println(refundResp.PayjsOrderID)
	fmt.Println(refundResp.OutTradeNo)
	fmt.Println(refundResp.TransactionID)
	fmt.Println(refundResp.Sign)
}
```

## Change logs

```
版本: v0.1
时间: 2020年1月27日
说明: 实现扫码支付, 查询订单, 关闭订单, 退款
```
