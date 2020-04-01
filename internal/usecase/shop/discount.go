package shop

import (
	"fmt"
)

const (
	StatusUnpublished = "未发布"
	StatusPublished   = "已发布"
)

type DiscountStrategy interface {
	CalculateDiscount() //计算方法
}

var discountStrategy map[string]DiscountStrategy

type Discount struct {
	SeqId        int64  `db:"seqId" json:"seqid"`
	Code         string `db:"code" json:"code"`                 //唯一编码
	Times        int    `db:"times" json:"times"`               //优惠次数（-10000为无限次）
	Priority     uint8  `db:"priority" json:"priority"`         //优惠优先级（数字越大等级越高）
	GoodsType    string `db:"goodstype" json:"goodstype"`       //优惠的商品类型
	GoodsCode    string `db:"goodscode" json:"goodscode"`       //优惠的商品编码
	StrategyKey  string `db:"strategykey" json:"strategykey"`   //优惠策略
	StrategyCode string `db:"strategycode" json:"strategycode"` //优惠策略的编码
	StartTime    string `db:"starttime" json:"starttime"`       //优惠开始时间
	EndTime      string `db:"endtime" json:"endtime"`           //优惠结束时间
	Status       string `db:"status" json:"status"`             //状态（未发布，已发布）
	Operator     string `db:"operator" json:"operator"`
	OperateTime  string `db:"operatetime" json:"operatetime"`
}

type DiscountOrderPayRel struct {
	SeqId          int64  `db:"seqId" json:"seqid"`
	DiscountCode   string `db:"discountcode" json:"discountcode"`     //优惠编码
	OrderGoodsCode string `db:"orderGoodscode" json:"orderGoodscode"` //支付订单编码
	Operator       string `db:"operator" json:"operator"`
	OperateTime    string `db:"operatetime" json:"operatetime"`
}

//注册优惠策略
func RegisterDiscountStrategy(key string, ds DiscountStrategy) error {
	if discountStrategy[key] != nil {
		return fmt.Errorf("key已被占用")
	} else {
		discountStrategy[key] = ds
		return nil
	}
}

func Add(times, priority int, goodsType, goodsCode, strategyKey, strategyCode, startTime, endTime string) {
	//TODO 状态设置为未发布
}

func Publish(seqId int64) {
	//TODO 状态设置为已发布
}

func Delete(seqId int64) {
	//TODO 只能删除未发布的优惠
}

func Update(seqId int64) {
	//TODO 只能更新未发布的优惠
}

//TODO 获取，根据页面具体设计
