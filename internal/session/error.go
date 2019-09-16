package session

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"runtime"

	"github.com/go-errors/errors"
)

// Error is a custom error
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

// AuthorizationError return 401 for unauthorized request
func AuthorizationError(ctx context.Context) Error {
	description := "Unauthorized, maybe invalid token."
	return createError(ctx, http.StatusAccepted, 401, description, nil)
}

// ForbiddenError return 401 for unauthorized request
func ForbiddenError(ctx context.Context) Error {
	description := http.StatusText(http.StatusForbidden)
	return createError(ctx, http.StatusAccepted, http.StatusForbidden, description, nil)
}

// BadRequestError means the request body is not a valid format.
func BadRequestError(ctx context.Context) Error {
	description := "The request body can’t be pasred as valid data."
	return createError(ctx, http.StatusAccepted, http.StatusBadRequest, description, nil)
}

// TransactionError means there is something wrong on database.
func TransactionError(ctx context.Context, err error) Error {
	description := http.StatusText(http.StatusInternalServerError)
	return createError(ctx, http.StatusInternalServerError, 10001, description, err)
}

// BadDataError means the request has invalid field.
func BadDataError(ctx context.Context) Error {
	description := "The request data has invalid field."
	return createError(ctx, http.StatusAccepted, 10002, description, nil)
}

// InvalidEmailFormatError means the email is invalid.
func InvalidEmailFormatError(ctx context.Context, email string) Error {
	description := fmt.Sprintf("Invalid email format %s.", email)
	return createError(ctx, http.StatusInternalServerError, 10010, description, nil)
}

// IdentityNonExistError means email or username is not existent.
func IdentityNonExistError(ctx context.Context) Error {
	description := "Email or Username is not exist."
	return createError(ctx, http.StatusAccepted, 10011, description, nil)
}

// InvalidPasswordError means the password is invalid.
func InvalidPasswordError(ctx context.Context) Error {
	description := "Password invalid."
	return createError(ctx, http.StatusAccepted, 10012, description, nil)
}

// PasswordTooSimpleError means the password is too simple.
func PasswordTooSimpleError(ctx context.Context) Error {
	description := "Password too simple, at least 8 characters required."
	return createError(ctx, http.StatusAccepted, 10013, description, nil)
}

// TooManyGroupsError means groups count is out of size.
func TooManyGroupsError(ctx context.Context) Error {
	description := "Too many groups."
	return createError(ctx, http.StatusAccepted, 10014, description, nil)
}

// InvalidImageDataError means image is invalid.
func InvalidImageDataError(ctx context.Context) Error {
	description := "Invalid image data."
	return createError(ctx, http.StatusAccepted, 10015, description, nil)
}

// TooManyGroupInvitationsError means group invitations count is out of size.
func TooManyGroupInvitationsError(ctx context.Context) Error {
	description := "Too Many Group Invitations."
	return createError(ctx, http.StatusAccepted, 10016, description, nil)
}

// InvalidGroupInvitationCodeError means group invitation code is invalid.
func InvalidGroupInvitationCodeError(ctx context.Context) Error {
	description := "Invalid Group Invitation Code."
	return createError(ctx, http.StatusAccepted, 10017, description, nil)
}

// RecaptchaVerifyError means recaptcha is invalid.
func RecaptchaVerifyError(ctx context.Context) Error {
	description := fmt.Sprintf("Recaptcha is invalid.")
	return createError(ctx, http.StatusAccepted, 10018, description, nil)
}

// VerificationCodeInvalidError means verification code is invalid
func VerificationCodeInvalidError(ctx context.Context) Error {
	description := "Invalid verification code."
	return createError(ctx, http.StatusAccepted, 10019, description, nil)
}

// ServerError means some server error are occurred.
func ServerError(ctx context.Context, err error) Error {
	description := http.StatusText(http.StatusInternalServerError)
	return createError(ctx, http.StatusInternalServerError, http.StatusInternalServerError, description, err)
}

// NotFoundError means resource is not found.
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
		trace = trace + "\n" + errors.Wrap(err, 1).ErrorStack()
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
