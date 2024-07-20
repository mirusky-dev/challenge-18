package core

import (
	"fmt"

	"github.com/gofiber/fiber/v2"
)

var ErrorHandler = func(ctx *fiber.Ctx, err error) error {
	if e, ok := err.(*Exception); ok {
		return ctx.Status(e.Status).JSON(e)
	} else if e, ok := err.(*fiber.Error); ok {
		return ctx.Status(e.Code).JSON(
			UserFriendlyException(
				WithCode("framework-error"),
				WithStatus(e.Code),
				WithMessage(e.Message),
			),
		)
	} else {
		return ctx.Status(500).JSON(
			Unexpected(
				WithError(err),
			),
		)
	}
}

type Severity string

const (
	Debug Severity = "debug"
	Info  Severity = "info"
	Warn  Severity = "warn"
	Error Severity = "error"
	Fatal Severity = "fatal"
)

type Exception struct {
	Status   int      `json:"status"`                       // HTTP Status Code
	Code     string   `json:"code"`                         // Application error code
	Message  string   `json:"message"`                      // User friendly message
	Err      string   `json:"error,omitempty" form:"error"` // Golang error
	Severity Severity `json:"severity"`                     // Exception level
}

func (e *Exception) Error() string {
	return e.Message
}

type UserFriendlyExceptionOption func(*Exception)

func WithStatus(status int) UserFriendlyExceptionOption {
	return func(h *Exception) {
		h.Status = status
	}
}

func WithCode(code string) UserFriendlyExceptionOption {
	return func(h *Exception) {
		h.Code = code
	}
}

func WithMessage(message string) UserFriendlyExceptionOption {
	return func(h *Exception) {
		h.Message = message
	}
}

func WithError(err error) UserFriendlyExceptionOption {
	return func(h *Exception) {
		if h.Err != "" {
			h.Err = fmt.Errorf("%s: %w", h.Err, err).Error()
		} else {
			h.Err = err.Error()
		}
	}
}

func WithSeverity(severity Severity) UserFriendlyExceptionOption {
	return func(h *Exception) {
		h.Severity = severity
	}
}

func NotAllowed(opts ...UserFriendlyExceptionOption) *Exception {
	defaultOpts := []UserFriendlyExceptionOption{
		WithCode("not-allowed"),
	}

	defaultOpts = append(defaultOpts, opts...)

	return NotImplemented(defaultOpts...)
}

func NotImplemented(opts ...UserFriendlyExceptionOption) *Exception {
	defaultOpts := []UserFriendlyExceptionOption{
		WithStatus(501),
		WithCode("not-implemented"),
		WithSeverity(Warn),
		WithMessage("The current method was not implemented"),
	}

	defaultOpts = append(defaultOpts, opts...)

	return UserFriendlyException(defaultOpts...)
}

func NotFound(opts ...UserFriendlyExceptionOption) *Exception {
	defaultOpts := []UserFriendlyExceptionOption{
		WithStatus(404),
		WithCode("not-found"),
		WithSeverity(Info),
		WithMessage("No entities found with given parameters"),
	}

	defaultOpts = append(defaultOpts, opts...)

	return UserFriendlyException(defaultOpts...)
}

func Forbidden(opts ...UserFriendlyExceptionOption) *Exception {
	defaultOpts := []UserFriendlyExceptionOption{
		WithStatus(403),
		WithCode("forbidden"),
		WithSeverity(Warn),
		WithMessage("You don't have permission for that!"),
	}

	defaultOpts = append(defaultOpts, opts...)

	return UserFriendlyException(defaultOpts...)
}

func MissingContext(opts ...UserFriendlyExceptionOption) *Exception {
	defaultOpts := []UserFriendlyExceptionOption{
		WithCode("missing-context"),
	}

	defaultOpts = append(defaultOpts, opts...)

	return UserFriendlyException(defaultOpts...)
}

func Unauthorized(opts ...UserFriendlyExceptionOption) *Exception {
	defaultOpts := []UserFriendlyExceptionOption{
		WithStatus(401),
		WithCode("unauthorized"),
		WithSeverity(Warn),
		WithMessage("You need to login first"),
	}

	defaultOpts = append(defaultOpts, opts...)

	return UserFriendlyException(defaultOpts...)
}

func BadRequest(opts ...UserFriendlyExceptionOption) *Exception {
	defaultOpts := []UserFriendlyExceptionOption{
		WithStatus(400),
		WithCode("bad-request"),
		WithSeverity(Info),
		WithMessage("Ops, something is wrong in the request"),
	}

	defaultOpts = append(defaultOpts, opts...)

	return UserFriendlyException(defaultOpts...)
}

func Unexpected(opts ...UserFriendlyExceptionOption) *Exception {
	defaultOpts := []UserFriendlyExceptionOption{
		WithStatus(500),
		WithCode("internal-server-error"),
		WithSeverity(Error),
		WithMessage("Something went wrong"),
	}

	defaultOpts = append(defaultOpts, opts...)

	return UserFriendlyException(defaultOpts...)
}

func UserFriendlyException(opts ...UserFriendlyExceptionOption) *Exception {
	const (
		defaultStatus   = 500
		defaultMessage  = "This is a friendly error, don't panic! Everything is under control"
		defaultCode     = "user-friendly-exception"
		defaultSeverity = Info
	)

	h := &Exception{
		Code:     defaultCode,
		Message:  defaultMessage,
		Status:   defaultStatus,
		Severity: defaultSeverity,
	}

	for _, opt := range opts {
		opt(h)
	}

	return h
}
