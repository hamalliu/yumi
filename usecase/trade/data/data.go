package data

import (
	"yumi/pkg/stores/dbc/mysqlx"
	"yumi/usecase/trade/service"
	"yumi/usecase/trade/thirdpf/alipay"
	"yumi/usecase/trade/thirdpf/wxpay"
)

// MysqlDB ...
type MysqlDB struct {
	curd mysqlx.CURD
	cli  *mysqlx.Client
}

var _ service.Data = &MysqlDB{}
var _ alipay.Data = &MysqlDB{}
var _ wxpay.Data = &MysqlDB{}

// New ...
func New(cli *mysqlx.Client) *MysqlDB {
	return &MysqlDB{curd: cli, cli: cli}
}
