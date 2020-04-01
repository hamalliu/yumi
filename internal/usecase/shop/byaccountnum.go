package shop

//已卖的产品
type SelledByUserNum struct {
	SeqId          int64  `db:"seqId" json:"seqid"`
	OrderGoodsCode string `db:"ordergoodscode" json:"ordergoodscode"`
}

type ByUserNumOrderGoodsRel struct {
	SeqId          int64  `db:"seqId" json:"seqid"`
	ByUserNumCode  string `db:"byusernumcode" json:"byusernumcode"`
	OrderGoodsCode string `db:"ordergoodscode" json:"ordergoodscode"`
}
