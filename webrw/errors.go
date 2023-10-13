package webrw

import (
	"errors"
	"net/http"
)

var (
	ErrInvalidRedirectStatusCode = errors.New("redirect status code is invalid")
)

var (
	ErrBadRequest                   = NewError(http.StatusBadRequest, http.StatusText(http.StatusBadRequest))                                     // 400
	ErrUnauthorized                 = NewError(http.StatusUnauthorized, http.StatusText(http.StatusUnauthorized))                                 // 401
	ErrPaymentRequired              = NewError(http.StatusPaymentRequired, http.StatusText(http.StatusPaymentRequired))                           // 402
	ErrForbidden                    = NewError(http.StatusForbidden, http.StatusText(http.StatusForbidden))                                       // 403
	ErrNotFound                     = NewError(http.StatusNotFound, http.StatusText(http.StatusNotFound))                                         // 404
	ErrMethodNotAllowed             = NewError(http.StatusMethodNotAllowed, http.StatusText(http.StatusMethodNotAllowed))                         // 405
	ErrNotAcceptable                = NewError(http.StatusNotAcceptable, http.StatusText(http.StatusNotAcceptable))                               // 406
	ErrProxyAuthRequired            = NewError(http.StatusProxyAuthRequired, http.StatusText(http.StatusProxyAuthRequired))                       // 407
	ErrRequestTimeout               = NewError(http.StatusRequestTimeout, http.StatusText(http.StatusRequestTimeout))                             // 408
	ErrConflict                     = NewError(http.StatusConflict, http.StatusText(http.StatusConflict))                                         // 409
	ErrGone                         = NewError(http.StatusGone, http.StatusText(http.StatusGone))                                                 // 410
	ErrLengthRequired               = NewError(http.StatusLengthRequired, http.StatusText(http.StatusLengthRequired))                             // 411
	ErrPreconditionFailed           = NewError(http.StatusPreconditionFailed, http.StatusText(http.StatusPreconditionFailed))                     // 412
	ErrRequestEntityTooLarge        = NewError(http.StatusRequestEntityTooLarge, http.StatusText(http.StatusRequestEntityTooLarge))               // 413
	ErrRequestURITooLong            = NewError(http.StatusRequestURITooLong, http.StatusText(http.StatusRequestURITooLong))                       // 414
	ErrUnsupportedMediaType         = NewError(http.StatusUnsupportedMediaType, http.StatusText(http.StatusUnsupportedMediaType))                 // 415
	ErrRequestedRangeNotSatisfiable = NewError(http.StatusRequestedRangeNotSatisfiable, http.StatusText(http.StatusRequestedRangeNotSatisfiable)) // 416
	ErrExpectationFailed            = NewError(http.StatusExpectationFailed, http.StatusText(http.StatusExpectationFailed))                       // 417
	ErrTeapot                       = NewError(http.StatusTeapot, http.StatusText(http.StatusTeapot))                                             // 418
	ErrMisdirectedRequest           = NewError(http.StatusMisdirectedRequest, http.StatusText(http.StatusMisdirectedRequest))                     // 421
	ErrUnprocessableEntity          = NewError(http.StatusUnprocessableEntity, http.StatusText(http.StatusUnprocessableEntity))                   // 422
	ErrLocked                       = NewError(http.StatusLocked, http.StatusText(http.StatusLocked))                                             // 423
	ErrFailedDependency             = NewError(http.StatusFailedDependency, http.StatusText(http.StatusFailedDependency))                         // 424
	ErrTooEarly                     = NewError(http.StatusTooEarly, http.StatusText(http.StatusTooEarly))                                         // 425
	ErrUpgradeRequired              = NewError(http.StatusUpgradeRequired, http.StatusText(http.StatusUpgradeRequired))                           // 426
	ErrPreconditionRequired         = NewError(http.StatusPreconditionRequired, http.StatusText(http.StatusPreconditionRequired))                 // 428
	ErrTooManyRequests              = NewError(http.StatusTooManyRequests, http.StatusText(http.StatusTooManyRequests))                           // 429
	ErrRequestHeaderFieldsTooLarge  = NewError(http.StatusRequestHeaderFieldsTooLarge, http.StatusText(http.StatusRequestHeaderFieldsTooLarge))   // 431
	ErrUnavailableForLegalReasons   = NewError(http.StatusUnavailableForLegalReasons, http.StatusText(http.StatusUnavailableForLegalReasons))     // 451

	ErrInternalServerError           = NewError(http.StatusInternalServerError, http.StatusText(http.StatusInternalServerError))                     // 500
	ErrNotImplemented                = NewError(http.StatusNotImplemented, http.StatusText(http.StatusNotImplemented))                               // 501
	ErrBadGateway                    = NewError(http.StatusBadGateway, http.StatusText(http.StatusBadGateway))                                       // 502
	ErrServiceUnavailable            = NewError(http.StatusServiceUnavailable, http.StatusText(http.StatusServiceUnavailable))                       // 503
	ErrGatewayTimeout                = NewError(http.StatusGatewayTimeout, http.StatusText(http.StatusGatewayTimeout))                               // 504
	ErrHTTPVersionNotSupported       = NewError(http.StatusHTTPVersionNotSupported, http.StatusText(http.StatusHTTPVersionNotSupported))             // 505
	ErrVariantAlsoNegotiates         = NewError(http.StatusVariantAlsoNegotiates, http.StatusText(http.StatusVariantAlsoNegotiates))                 // 506
	ErrInsufficientStorage           = NewError(http.StatusInsufficientStorage, http.StatusText(http.StatusInsufficientStorage))                     // 507
	ErrLoopDetected                  = NewError(http.StatusLoopDetected, http.StatusText(http.StatusLoopDetected))                                   // 508
	ErrNotExtended                   = NewError(http.StatusNotExtended, http.StatusText(http.StatusNotExtended))                                     // 510
	ErrNetworkAuthenticationRequired = NewError(http.StatusNetworkAuthenticationRequired, http.StatusText(http.StatusNetworkAuthenticationRequired)) // 511
)

type Error struct {
	Code    int    `json:"-"`
	Message string `json:"message"`
}

func NewError(status int, msg string) *Error {
	return &Error{
		Code:    status,
		Message: msg,
	}
}

func (e *Error) Error() string {
	return e.Message
}
