package errs

import "encoding/json"

type Error struct {
	Code          string      `json:"code"`
	Message       string      `json:"message"`
	InternalError interface{} `json:"-"`
}

func (e *Error) Error() string {
	if e == nil {
		return ""
	}
	indent, err := json.MarshalIndent(e, " ", "  ")
	if err != nil {
		return ""
	}

	return string(indent)
}

var (
	BadRequest = &Error{
		Code:    "40001",
		Message: "400 wrong form",
	}
	SendSmsPleaseWait = &Error{
		Code:    "40002",
		Message: "SMS has been sent, please wait",
	}
	CaptchaCode = &Error{
		Code:    "4008",
		Message: "Incorrect verification code or expired",
	}
	CaptchaCode2 = &Error{
		Code:    "4009",
		Message: "Verification code error Or Expired",
	}
	PleaseSignIn = &Error{
		Code:    "401",
		Message: "401 please sign in",
	}
	NotData = &Error{
		Code:    "404",
		Message: "Not Data",
	}
	LoginFailed = &Error{
		Code:    "41001",
		Message: "Login failed Incorrect password or User does not exist",
	}
	WithdrawalPasswordWrong = &Error{
		Code:    "41003",
		Message: "Withdrawal password wrong",
	}
	ExUser = &Error{
		Code:    "41002",
		Message: "User already exists",
	}
	SqlSystemError = func(err interface{}) *Error {
		return &Error{
			Code:          "5002",
			Message:       "System upgrade, please try again later",
			InternalError: err,
		}
	}
	ReptileError = &Error{
		Code:    "5002",
		Message: "In order to record your crimes, our company will prosecute you for illegal intrusion into the computer information system",
	}
	SystemError = func(err interface{}) *Error {
		return &Error{
			Code:          "5001",
			Message:       "System upgrade, please try again later",
			InternalError: err,
		}
	}
)

func NewError(code string, message string, body interface{}) *Error {
	return &Error{
		Code:    code,
		Message: message,
	}
}
