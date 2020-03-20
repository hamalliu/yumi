package goods

type GoodsInfo interface {
	Produce() //
	Destroy()
}

type Goods struct {
	SeqId          int64  `db:"seqId" json:"seqid"`
	Code           string `db:"code" json:"code"` //
	GoodsKey       string `db:"goodskey" json:"goodskey"`
	GoodsCode      string `db:"goodscode" json:"goodscode"`
	ProductionDate string `db:"productiondate" json:"productiondate"` //生产日期
	ExpirationDate string `db:"expirationdate" json:"expirationdate"` //过期日期
}
