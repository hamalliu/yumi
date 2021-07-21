package usecase

import (
	"yumi/pkg/stores/dbc/mysqlx"
	"yumi/pkg/stores/mgoc"
	"yumi/usecase/trade"
	tradeservice "yumi/usecase/trade/service"
	"yumi/usecase/trade/thirdpf/wxpay"
	"yumi/usecase/user"
	"yumi/usecase/media"
	userservice "yumi/usecase/user/service"
	mediaservice "yumi/usecase/media/service"
)

var tradeSrv *tradeservice.Service

func InstallTrade(mongoC *mgoc.Client, mysqlC *mysqlx.Client, wxMwebConf wxpay.MwebConfig) {
	var err error
	tradeSrv, err = trade.Usecase(mongoC, mysqlC, wxMwebConf)
	if err != nil {
		panic(err)
	}
}

func Trade() *tradeservice.Service {
	return tradeSrv
}

var userSrv *userservice.Service

func InstallUser(mongoC *mgoc.Client) {
	var err error
	userSrv, err = user.Usecase(mongoC)
	if err != nil {
		panic(err)
	}
}

func User() *userservice.Service {
	return userSrv
}

var mediaSrv *mediaservice.Service

func InstallMedia(mysqlC *mysqlx.Client) {
	var err error
	mediaSrv, err = media.Usecase(mysqlC)
	if err != nil {
		panic(err)
	}
}

func Media() *mediaservice.Service {
	return mediaSrv
}
