package status

type messageEntry int

const (
	// =====================user=====================//
	// 密码错误
	PasswordIncorrect messageEntry = iota //start
	// 用户名格式不正确
	UserFmtIncorrect
	// 密码格式不正确
	PasswordFmtIncorrect
	// 用户名不存在
	UserNotFound
	//================================================//


	// =====================trade=====================//
	// 订单超时
	OrderTimeout
	// 订单已存在
	OrderAlreadyExists
	// 订单已完成，不能取消
	OrderFinishedRefuseCancel
	//================================================//

	
	// =====================media=====================//
	// 文件太大
	FileIsTooLarge
	//================================================//
)
