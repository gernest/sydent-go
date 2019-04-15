package models

import (
	"encoding/json"
	"errors"
)

var ErrIncorrectClientSecret = errors.New("db: incorrect client secret")
var ErrSessionExpired = errors.New("db: session expired")
var ErrSessionNotValidated = errors.New("db: session not validated")
var ErrIncorrectToken = errors.New("bad token")
var ErrNoSignature = errors.New("no signatures found")
var ErrNoMatchingSignature = errors.New("no matching signatures found")

// matrix error codes
const (
	ErrUnrecognized                = "M_UNRECOGNIZED"
	ErrUnauthorized                = "M_UNAUTHORIZED"
	ErrForbidden                   = "M_FORBIDDEN"
	ErrBadJSON                     = "M_BAD_JSON"
	ErrNotJSON                     = "M_NOT_JSON"
	ErrUserInUse                   = "M_USER_IN_USE"
	ErrRoomInUse                   = "M_ROOM_IN_USE"
	ErrBadPagination               = "M_BAD_PAGINATION"
	ErrBadState                    = "M_BAD_STATE"
	ErrUnknown                     = "M_UNKNOWN"
	ErrNotFound                    = "M_NOT_FOUND"
	ErrMissingToken                = "M_MISSING_TOKEN"
	ErrUnknownToken                = "M_UNKNOWN_TOKEN"
	ErrGuestAccessForbidden        = "M_GUEST_ACCESS_FORBIDDEN"
	ErrLimitExceeded               = "M_LIMIT_EXCEEDED"
	ErrCaptchaNeeded               = "M_CAPTCHA_NEEDED"
	ErrCaptchaInvalid              = "M_CAPTCHA_INVALID"
	ErrMissingParam                = "M_MISSING_PARAM"
	ErrInvalidParam                = "M_INVALID_PARAM"
	ErrTooLarge                    = "M_TOO_LARGE"
	ErrExclusive                   = "M_EXCLUSIVE"
	ErrThreepidAuthFailed          = "M_THREEPID_AUTH_FAILED"
	ErrThreepidInUse               = "M_THREEPID_IN_USE"
	ErrThreepidNotFound            = "M_THREEPID_NOT_FOUND"
	ErrThreepidDenied              = "M_THREEPID_DENIED"
	ErrInvalidUsername             = "M_INVALID_USERNAME"
	ErrServerNotTrusted            = "M_SERVER_NOT_TRUSTED"
	ErrConsentNotGiven             = "M_CONSENT_NOT_GIVEN"
	ErrCannotLeaveServerNoticeRoom = "M_CANNOT_LEAVE_SERVER_NOTICE_ROOM"
	ErrResourceLimitExceeded       = "M_RESOURCE_LIMIT_EXCEEDED"
	ErrUnsupportedRoomVersion      = "M_UNSUPPORTED_ROOM_VERSION"
	ErrIncompatibleRoomVersion     = "M_INCOMPATIBLE_ROOM_VERSION"
	ErrWrongRoomKeysVersion        = "M_WRONG_ROOM_KEYS_VERSION"
	ErrNoValidSession              = "M_NO_VALID_SESSION"
	ErrSessionExpiredCode          = "M_SESSION_EXPIRED"
	ErrSessionNotValidatedCode     = "M_SESSION_NOT_VALIDATED"
	ErrUnknownPeer                 = "M_UNKNOWN_PEER"
	ErrVerificationFailed          = "M_VERIFICATION_FAILED"
	ErrInvalidEmail                = "M_INVALID_EMAIL"
	ErrEmailSendError              = "M_EMAIL_SEND_ERROR"
	ErrIncorrectClientSecretCode   = "M_INCORRECT_CLIENT_SECRET"
)

// Error implements error interface, that wraps around the response from the
// matrix server.
type Error struct {
	Code string `json:"errcode"`
	Err  string `json:"error"`
}

// NewError returns Error with errcode and error fields set to code and err
// respectively.
func NewError(code, err string) Error {
	return Error{
		Code: code,
		Err:  err,
	}
}

func (e Error) Error() string {
	b, _ := json.Marshal(e)
	return string(b)
}
