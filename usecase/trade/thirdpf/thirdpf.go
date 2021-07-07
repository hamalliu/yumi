package thirdpf

import (
	"errors"
	
	"yumi/pkg/status"
	"yumi/usecase/trade/entity"
)

const (
	//TradeWayAliPayPage ...
	TradeWayAliPayPage = Way("alipay_page")

	//TradeWayWxPayAPP ...
	TradeWayWxPayAPP = Way("wxpay_app")
	//TradeWayWxPayJSAPI ...
	TradeWayWxPayJSAPI = Way("wxpay_jsapi")
	//TradeWayWxPayMWEB ...
	TradeWayWxPayMWEB = Way("wxpay_mweb")
	//TradeWayWxPayNATIVE1 ...
	TradeWayWxPayNATIVE1 = Way("wxpay_native1")
	//TradeWayWxPayNATIVE2 ...
	TradeWayWxPayNATIVE2 = Way("wxpay_native2")
)

//Way ...
type Way string

type Trades struct {
	trades map[Way]entity.Trade
}

// NewThirdpf ...
func New() *Trades {
	t := &Trades{}
	t.trades = make(map[Way]entity.Trade)
	return t
}

// AddTrade 新增交易方式
func (t *Trades) AddTrade(way Way, trade entity.Trade) {
	t.trades[way] = trade
}

// GetTrade 获取交易方式
func (t *Trades) GetTrade(way Way) (entity.Trade, error) {
	if t.trades[way] == nil {
		return nil, status.Internal().WithDetails(errors.New("unsupported trade way"))
	}
	return t.trades[way], nil
}
