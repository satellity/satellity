package session

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"
)

type Error struct {
	Status      int    `json:"status"`
	Code        int    `json:"code"`
	Description string `json:"description"`
	trace       string
}

func (sessionError Error) Error() string {
	str, err := json.Marshal(sessionError)
	if err != nil {
		log.Panicln(err)
	}
	return string(str)
}

func BadRequestError(ctx context.Context) Error {
	description := "The request body can’t be pasred as valid data."
	return createError(ctx, http.StatusAccepted, http.StatusBadRequest, description, nil)
}

func TransactionError(ctx context.Context, err error) Error {
	description := http.StatusText(http.StatusInternalServerError)
	return createError(ctx, http.StatusInternalServerError, 10001, description, err)
}

func BadDataError(ctx context.Context) Error {
	description := "The request data has invalid field."
	return createError(ctx, http.StatusAccepted, 10002, description, nil)
}

func InvalidEmailFormatError(ctx context.Context, email string) Error {
	description := fmt.Sprintf("Invalid email format %s.", email)
	return createError(ctx, http.StatusInternalServerError, 10010, description, nil)
}

func IdentityNonExistError(ctx context.Context) Error {
	description := "Email or Username is not exist."
	return createError(ctx, http.StatusAccepted, 10011, description, nil)
}

func InvalidPasswordError(ctx context.Context) Error {
	description := "Password invalid."
	return createError(ctx, http.StatusAccepted, 10012, description, nil)
}

func PasswordTooSimpleError(ctx context.Context) Error {
	description := "Password too simple, at least 8 characters required."
	return createError(ctx, http.StatusAccepted, 10013, description, nil)
}

func ServerError(ctx context.Context, err error) Error {
	description := http.StatusText(http.StatusInternalServerError)
	return createError(ctx, http.StatusInternalServerError, http.StatusInternalServerError, description, err)
}

func NotFoundError(ctx context.Context) Error {
	description := http.StatusText(http.StatusNotFound)
	return createError(ctx, http.StatusAccepted, http.StatusNotFound, description, nil)
}

func createError(ctx context.Context, status, code int, description string, err error) Error {
	pc, file, line, _ := runtime.Caller(2)
	funcName := runtime.FuncForPC(pc).Name()
	trace := fmt.Sprintf("[ERROR %d] %s\n%s:%d %s", code, description, file, line, funcName)
	if err != nil {
		if sessionError, ok := err.(Error); ok {
			trace = trace + "\n" + sessionError.trace
		} else {
			trace = trace + "\n" + err.Error()
		}
	}

	if ctx != nil {
		if logger := Logger(ctx); logger != nil {
			logger.Error(trace)
		}
	}

	return Error{
		Status:      status,
		Code:        code,
		Description: description,
		trace:       trace,
	}
}
