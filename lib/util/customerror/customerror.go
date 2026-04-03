package customerror

import (
	"encoding/json"
	"errors"
	"fmt"
	"log"
	"log/slog"
)

const (
	featDigit = 2
	caseDigit = 3
)

const (
	LogErrorKey = "error"
	LogPanicKey = "panic"
)

type Model struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (m Model) Error() string {
	jsonData, _ := json.Marshal(m)
	jsonDataStr := string(jsonData)
	return jsonDataStr
}

func UnmarshalError(werr error) *Model {
	var customErr Model

	uerr := errors.Unwrap(errors.Unwrap(werr))
	if uerr == nil {
		log.Println("Cannot unwrap error:", werr)
		return &customErr
	}

	err := json.Unmarshal([]byte(uerr.Error()), &customErr)
	if err != nil {
		log.Println("Cannot unmarshal the error: ", err)
		return &customErr
	}

	return &customErr
}

func New(featNum int, caseNum int, msg string) Model {
	return Model{
		Code:    fmt.Sprintf("%0*d%0*d", featDigit, featNum, caseDigit, caseNum),
		Message: msg,
	}
}

// NewInternalErr creates a temporary error with code "00000" for development/debugging
// This is for internal communication between developers and will not be displayed to users
// Use this when you need to quickly communicate an error without creating a formal error code
func NewInternalErr(msg string) Model {
	return Model{
		Code:    "00000",
		Message: msg,
	}
}

func Log(logger *slog.Logger, customError Model, realError error) Model {
	if realError != nil {
		logger.Error(customError.Message, LogErrorKey, realError)
	} else {
		logger.Error(customError.Message)
	}

	return customError
}
