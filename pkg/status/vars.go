package status

import "yumi/pkg/codes"

var (
	m = NewMessages()

	// OperationSuccess 操作成功
	OKMessage                 = m.NewMessageID("operation success", "操作成功", "操作成功")
	CanceledMessage           = m.NewMessageID("operation canceled", "操作被取消", "操作被取消")
	UnknownMessage            = m.NewMessageID("an unknown error", "未知错误", "未知錯誤")
	InvalidArgumentMessage    = m.NewMessageID("invalid argument", "参数错误", "參數錯誤")
	DeadlineExceededMessage   = m.NewMessageID("processing timed out", "处理已超时", "處理已超時")
	NotFoundMessage           = m.NewMessageID("not found data", "未找到数据", "未找到數據")
	AlreadyExistsMessage      = m.NewMessageID("data already exists", "数据已存在", "數據已存在")
	PermissionDeniedMessage   = m.NewMessageID("permission denied", "没有权限", "沒有權限")
	ResourceExhaustedMessage  = m.NewMessageID("resource exhausted", "服务器资源耗尽", "服務器資源耗盡")
	FailedPreconditionMessage = m.NewMessageID("failed precondition", "前置条件不满足", "前置條件不滿足")
	AbortedMessage            = m.NewMessageID("operation aborted", "操作被中断", "操作被中斷")
	OutOfRangeMessage         = m.NewMessageID("out of range", "数组越界", "數組越界")
	UnimplementedMessage      = m.NewMessageID("operation unimplemented", "操作未声明", "操作未聲明")
	InternalMessage           = m.NewMessageID("server internal error", "服务器内部错误", "服務器内部錯誤")
	UnavailableMessage        = m.NewMessageID("service unavailable", "服务不可用", "服務不可用")
	DataLossMessage           = m.NewMessageID("data loss", "数据已损坏或丢失", "數據已損壞或丟失")
	UnauthenticatedMessage    = m.NewMessageID("operation unauthenticated", "操作未认证", "操作未認證")
)

func init() {
	m.InitI18N()
}

var (
	// OK ...
	OK = func() *Status { return new(codes.OK, OKMessage) }
	// Canceled ...
	Canceled = func() *Status { return new(codes.Canceled, CanceledMessage) }
	// Unknown ...
	Unknown = func() *Status { return new(codes.Unknown, UnknownMessage) }
	// InvalidArgument ...
	InvalidArgument = func() *Status { return new(codes.InvalidArgument, InvalidArgumentMessage) }
	// DeadlineExceeded ...
	DeadlineExceeded = func() *Status { return new(codes.DeadlineExceeded, DeadlineExceededMessage) }
	// NotFound ...
	NotFound = func() *Status { return new(codes.NotFound, NotFoundMessage) }
	// AlreadyExists ...
	AlreadyExists = func() *Status { return new(codes.AlreadyExists, AlreadyExistsMessage) }
	// PermissionDenied ...
	PermissionDenied = func() *Status { return new(codes.PermissionDenied, PermissionDeniedMessage) }
	// ResourceExhausted ...
	ResourceExhausted = func() *Status { return new(codes.ResourceExhausted, ResourceExhaustedMessage) }
	// FailedPrecondition ...
	FailedPrecondition = func() *Status { return new(codes.FailedPrecondition, FailedPreconditionMessage) }
	// Aborted ...
	Aborted = func() *Status { return new(codes.Aborted, AbortedMessage) }
	// OutOfRange ...
	OutOfRange = func() *Status { return new(codes.OutOfRange, OutOfRangeMessage) }
	// Unimplemented ...
	Unimplemented = func() *Status { return new(codes.Unimplemented, UnimplementedMessage) }
	// Internal ...
	Internal = func() *Status { return new(codes.Internal, InternalMessage) }
	// Unavailable ...
	Unavailable = func() *Status { return new(codes.Unavailable, UnavailableMessage) }
	// DataLoss ...
	DataLoss = func() *Status { return new(codes.DataLoss, DataLossMessage) }
	// Unauthenticated ...
	Unauthenticated = func() *Status { return new(codes.Unauthenticated, UnauthenticatedMessage) }
)
