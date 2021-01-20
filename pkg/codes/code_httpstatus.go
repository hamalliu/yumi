package codes

import "net/http"

// HTTPStatus 返回http状态码
func (c Code) HTTPStatus() int {
	switch c {
	case OK:
		return http.StatusOK
	case Canceled:
		return 499
	case Unknown:
		return http.StatusInternalServerError
	case InvalidArgument:
		return http.StatusBadRequest
	case DeadlineExceeded:
		return http.StatusGatewayTimeout
	case NotFound:
		return http.StatusNotFound
	case AlreadyExists:
		return http.StatusConflict
	case PermissionDenied:
		return http.StatusForbidden
	case ResourceExhausted:
		return http.StatusTooManyRequests
	case FailedPrecondition:
		return http.StatusPreconditionFailed
	case Aborted:
		return http.StatusConflict
	case OutOfRange:
		return http.StatusBadRequest
	case Unimplemented:
		return http.StatusNotImplemented
	case Internal:
		return http.StatusInternalServerError
	case Unavailable:
		return http.StatusServiceUnavailable
	case DataLoss:
		return http.StatusInternalServerError
	case Unauthenticated:
		return http.StatusUnauthorized
	default:
		return http.StatusInternalServerError
	}
}