package status

import "yumi/pkg/codes"

var (
	// OK ...
	OK = func() *Status { return new(codes.OK, codes.OK.String()) }
	// Canceled ...
	Canceled = func() *Status { return new(codes.Canceled, codes.Canceled.String()) }
	// Unknown ...
	Unknown = func() *Status { return new(codes.Unknown, codes.Unknown.String()) }
	// InvalidArgument ...
	InvalidArgument = func() *Status { return new(codes.InvalidArgument, codes.InvalidArgument.String()) }
	// DeadlineExceeded ...
	DeadlineExceeded = func() *Status { return new(codes.DeadlineExceeded, codes.DeadlineExceeded.String()) }
	// NotFound ...
	NotFound = func() *Status { return new(codes.NotFound, codes.NotFound.String()) }
	// AlreadyExists ...
	AlreadyExists = func() *Status { return new(codes.AlreadyExists, codes.AlreadyExists.String()) }
	// PermissionDenied ...
	PermissionDenied = func() *Status { return new(codes.PermissionDenied, codes.PermissionDenied.String()) }
	// ResourceExhausted ...
	ResourceExhausted = func() *Status { return new(codes.ResourceExhausted, codes.ResourceExhausted.String()) }
	// FailedPrecondition ...
	FailedPrecondition = func() *Status { return new(codes.FailedPrecondition, codes.FailedPrecondition.String()) }
	// Aborted ...
	Aborted = func() *Status { return new(codes.Aborted, codes.Aborted.String()) }
	// OutOfRange ...
	OutOfRange = func() *Status { return new(codes.OutOfRange, codes.OutOfRange.String()) }
	// Unimplemented ...
	Unimplemented = func() *Status { return new(codes.Unimplemented, codes.Unimplemented.String()) }
	// Internal ...
	Internal = func() *Status { return new(codes.Internal, codes.Internal.String()) }
	// Unavailable ...
	Unavailable = func() *Status { return new(codes.Unavailable, codes.Unavailable.String()) }
	// DataLoss ...
	DataLoss = func() *Status { return new(codes.DataLoss, codes.DataLoss.String()) }
	// Unauthenticated ...
	Unauthenticated = func() *Status { return new(codes.Unauthenticated, codes.Unauthenticated.String()) }
	// InvalidRequest ...
	InvalidRequest = func() *Status { return new(codes.InvalidRequest, codes.InvalidRequest.String()) }
)
