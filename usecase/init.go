package usecase

import (
	trade_db "yumi/usecase/trade/db"
	trade_platform "yumi/usecase/trade/platform"
)

// Init trade db
func Init() {
	trade_db.Init()
	trade_platform.Init()
}
