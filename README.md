# Payjs SDK for Go

## 当前实现

- [x] 扫码支付
- [x] 付款码支付
- [x] 查询订单
- [x] 关闭订单
- [x] 撤销订单
- [x] 退款
- [ ] 人脸支付
- [ ] 分账支付
- [ ] JSAPI支付
- [ ] 收银台支付

## 获取模块

```shell
> go get github.com/markity/payjs
```

## 基本使用

- [扫码支付(native)](#扫码支付native)
- [付款码支付(micropay)](#付款码支付micropay)
- [查询订单(check)](#查询订单check)
- [关闭订单(close)](#关闭订单close)
- [撤销订单(reverse)](#撤销订单reverse)
- [退款(refund)](#退款refund)

### 扫码支付(native)

**使用场景: 网站生成二维码, 用户扫码支付后回调, 网站执行相关逻辑**

```go
package main

import (
	"fmt"
	"github.com/markity/payjs"
)

func main() {
	mch := payjs.NewMch("your-mchid", "your-mchkey")
	// 说明: ReturnCode不为1 或 签名错误均视为error, 此模块所有方法均自动校验签名
	nativePayResp, err := mch.NativePay(payjs.NativePayInfo{TotalFee: 10, OutTradeNo: "2020_1_27_001", Body: "支付测试"})
	if err != nil {
		fmt.Printf("失败: %v\n", err)
		return
	}

	fmt.Println(nativePayResp.ReturnCode)
	fmt.Println(nativePayResp.ReturnMsg)
	fmt.Println(nativePayResp.PayjsOrderID)
	fmt.Println(nativePayResp.Qrcode)
	fmt.Println(nativePayResp.CodeUrl)
	fmt.Println(nativePayResp.Sign)
}
```

### 付款码支付(micropay)

**使用场景: 付款者打开微信首付款, 商家使用扫码枪或其他设备扫描用户auth_code, 程序拿到auth_code可调起支付**

```go
package main

import (
	"fmt"
	"github.com/markity/payjs"
)

func main() {
	mch := payjs.NewMch("your-mchid", "your-mchkey")
	microPayResp, err := mch.MicroPay(payjs.MicroPayInfo{TotalFee: 10, OutTradeNo: "2020_1_28_001", AuthCode: "付款者微信收付款的auth-code"})
	if err != nil {
		fmt.Printf("失败: %v\n", err)
		return
	}

	fmt.Println(microPayResp.ReturnCode)
	fmt.Println(microPayResp.ReturnMsg)
	fmt.Println(microPayResp.PayjsOrderID)
	fmt.Println(microPayResp.Sign)
}
```

> 注意1: 付款码支付不支持回调, 需自行调用CheckOrder检查是否完成支付

> 注意2: 不论请求成功失败, microPayResp.ReturnCode恒为0, 成功与否取决于microPayResp.ReturnMsg是否为`需要用户输入支付密码`

### 查询订单(check)

**使用场景: 用户发起支付后，可通过本接口发起订单查询来确认订单状态**

```go
package main

import (
	"fmt"
	"github.com/markity/payjs"
)

func main() {
	mch := payjs.NewMch("your-mchid", "your-mchkey")
	checkResp, err := mch.CheckOrder(payjs.CheckOrderInfo{PayjsOrderID: "order-id"})
	if err != nil {
		fmt.Printf("失败: %v\n", err)
		return
	}

	fmt.Println(checkResp.ReturnCode)
	fmt.Println(checkResp.OutTradeNo)
	fmt.Println(checkResp.TransactionID)
	fmt.Println(checkResp.Status)
	fmt.Println(checkResp.Openid)
	fmt.Println(checkResp.TotalFee)
	fmt.Println(checkResp.PaidTime)
	fmt.Println(checkResp.Attach)
	fmt.Println(checkResp.Sign)
}
```

> 注意: 一个支付成功的订单, 如果接收到check查询, 代表已经信息成功送达.此时如果异步通知尚未执行或正在执行, 则会中断执行, 此订单不再回调

### 关闭订单(close)

```go
package main

import (
	"fmt"
	"github.com/markity/payjs"
)

func main() {
	mch := payjs.NewMch("your-mchid", "your-mchkey")
	closeResp, err := mch.CloseOrder(payjs.CloseOrderInfo{PayjsOrderID: "order-id"})
	if err != nil {
		fmt.Printf("失败: %v\n", err)
		return
	}

	fmt.Println(closeResp.ReturnCode)
	fmt.Println(closeResp.ReturnMsg)
	fmt.Println(closeResp.Sign)
}
```

> 提醒: 仅能关闭30天以内的未支付订单, 订单发起后, 并非一定要关闭

### 撤销订单(reverse)

**使用场景: 针对人脸支付和付款码支付的异常订单(例如无法查询或确定订单状态), 对订单进行撤销, 若订单已经支付, 则自动原路还款**

```go
package main

import (
	"fmt"
	"github.com/markity/payjs"
)

func main() {
	mch := payjs.NewMch("your-mchid", "your-mchkey")
	reverseRsep, err := mch.ReverseOrder(payjs.ReverseOrderInfo{PayjsOrderID: "order-id"})
	if err != nil {
		fmt.Printf("失败: %v\n", err)
		return
	}

	fmt.Println(reverseRsep.ReturnCode)
	fmt.Println(reverseRsep.ReturnMsg)
	fmt.Println(reverseRsep.ReturnCode)

}
```

> 注意: 撤销订单无法适用于正常订单, 只能应用于7天以内的异常订单: 何时使用?参见https://pay.weixin.qq.com/wiki/doc/api/micropay.php?chapter=5_4

### 退款(refund)

**使用场景: 对于已经付款的订单, 可以退款, 资金自动原路返回**

```go
package main

import (
	"fmt"
	"github.com/markity/payjs"
)

func main() {
	mch := payjs.NewMch("your-mchid", "your-mchkey")
	refundResp, err := mch.Refund(payjs.RefundInfo{PayjsOrderID: "order-id"})
	if err != nil {
		fmt.Printf("失败: %v\n", err)
		return
	}

	fmt.Println(refundResp.ReturnCode)
	fmt.Println(refundResp.ReturnMsg)
	fmt.Println(refundResp.Sign)
}
```

> 注意: 若向已经退款的订单重复退款, return_code为1, 但无sign字段

## Change logs

```
版本: v0.3.1
时间: 2020年2月7日
内容:
  Native统一改名NativePay
  修正README中的错误
```

```
版本: v0.3
时间: 2020年1月31日
内容:
  移除丢弃响应中的不必要字段(原样返回的字段), 保留必要字段, 防止payjs接口频繁更新, 某些字段消失, 便于维护
  移除tool.go中打印出的调试信息
```

```
版本: v0.2.2
时间: 2020年1月30日
内容:
  优化签名校验算法, 可以适应于更为复杂的数据结构
  payjs.DEBUG由常量改为变量
```

```
版本: v0.2.1
时间: 2020年1月29日
内容:
  mch.Native改名mch.NativePay
  新增payjs.DEBUG配置, 开启后发送的json数据以及服务器响应的json数据将被输出
```

```
版本: v0.2
时间: 2020年1月28日
内容:
  修复订单重复退款响应的签名无法正常校验的bug
  实现了撤销订单, 付款码支付功能
```

```
版本: v0.1
时间: 2020年1月27日
内容:
  实现扫码支付, 查询订单, 关闭订单, 退款
```
