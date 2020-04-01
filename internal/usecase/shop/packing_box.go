package shop

type PackingBox struct {
	SeqId       int64  `db:"seqId" json:"seqid"`
	GoodsCode   string `db:"goodscode" json:"goodscode"`     //商品编码
	Price       int    `db:"price" json:"price"`             //单价
	Body        string `db:"body" json:"body"`               //描述
	Detail      string `db:"detail" json:"detail"`           //详情
	Operator    string `db:"operator" json:"operator"`       //操作人
	OperateTime string `db:"operatetime" json:"operatetime"` //操作时间
}
