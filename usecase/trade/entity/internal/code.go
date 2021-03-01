package internal

import (
	"fmt"
	"strings"
	"sync/atomic"
	"time"

	"yumi/pkg/random"
)

// GetOutTradeNo 获取商户订单号
func GetOutTradeNo() string {
	prefix := strings.ReplaceAll(time.Now().Format("06121545.999999999"), ".", "")
	return fmt.Sprintf("%s%s", prefix, random.Get(7, random.NUMBER))
}

//CodeType 生成订单号
type CodeType uint8

//编码类型
const (
	//OrderPayCode 支付订单号
	OrderPayCode CodeType = iota
	//OrderRefundCode 退款订单号
	OrderRefundCode
)

var count uint64

// GetCode 获取订单号
func GetCode(codeType CodeType) string {
	prefix := strings.ReplaceAll(time.Now().Format("06121545.999"), ".", "")
	random := random.Get(3, random.NUMBER)
	if count >= 100 {
		count = 0
	}

	atomic.AddUint64(&count, 1)
	return fmt.Sprintf("%s%d%d%s", prefix, codeType, count, random)
}
