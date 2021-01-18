package status

import "yumi/pkg/codes"

var (
	// OK ...
	OK = func() *Status { return New(codes.OK, codes.OK.String()) }
	// Canceled ...
	Canceled = func() *Status { return New(codes.Canceled, codes.Canceled.String()) }
	// Unknown ...
	Unknown = func() *Status { return New(codes.Unknown, codes.Unknown.String()) }
	// InvalidArgument ...
	InvalidArgument = func() *Status { return New(codes.InvalidArgument, codes.InvalidArgument.String()) }
	// DeadlineExceeded ...
	DeadlineExceeded = func() *Status { return New(codes.DeadlineExceeded, codes.DeadlineExceeded.String()) }
	// NotFound ...
	NotFound = func() *Status { return New(codes.NotFound, codes.NotFound.String()) }
	// AlreadyExists ...
	AlreadyExists = func() *Status { return New(codes.AlreadyExists, codes.AlreadyExists.String()) }
	// PermissionDenied ...
	PermissionDenied = func() *Status { return New(codes.PermissionDenied, codes.PermissionDenied.String()) }
	// ResourceExhausted ...
	ResourceExhausted = func() *Status { return New(codes.ResourceExhausted, codes.ResourceExhausted.String()) }
	// FailedPrecondition ...
	FailedPrecondition = func() *Status { return New(codes.FailedPrecondition, codes.FailedPrecondition.String()) }
	// Aborted ...
	Aborted = func() *Status { return New(codes.Aborted, codes.Aborted.String()) }
	// OutOfRange ...
	OutOfRange = func() *Status { return New(codes.OutOfRange, codes.OutOfRange.String()) }
	// Unimplemented ...
	Unimplemented = func() *Status { return New(codes.Unimplemented, codes.Unimplemented.String()) }
	// Internal ...
	Internal = func() *Status { return New(codes.Internal, codes.Internal.String()) }
	// Unavailable ...
	Unavailable = func() *Status { return New(codes.Unavailable, codes.Unavailable.String()) }
	// DataLoss ...
	DataLoss = func() *Status { return New(codes.DataLoss, codes.DataLoss.String()) }
	// Unauthenticated ...
	Unauthenticated = func() *Status { return New(codes.Unauthenticated, codes.Unauthenticated.String()) }
)