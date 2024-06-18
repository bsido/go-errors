package errors

import (
	"errors"
	"fmt"
	"log"
	"maps"
	"strings"
)

type Error struct {
	message string
	cause   error
	code    string
	helps   []string
	notes   []string

	additionalTemplateData map[string]any

	wrapped []error

	init *Init
}

func New(message string) *Error {
	return &Error{
		message: message,
		wrapped: make([]error, 0),
	}
}

func Newf(format string, args ...any) *Error {
	return New(fmt.Sprintf(format, args...))
}

func Extend(err error) *Error {
	var e *Error
	if errors.As(err, &e) {
		return e
	}

	return New(err.Error())
}

func ExtendWithMessage(err error, message string) *Error {
	var e *Error
	if errors.As(err, &e) {
		// override the original message
		e.message = message
		return e
	}

	return New(err.Error())
}

func (e *Error) initializer(b *Init) *Error {
	e.init = b

	return e
}

func (e *Error) Wrap(err error) *Error {
	e.wrapped = append(e.wrapped, err)

	return e
}

func (e *Error) Cause(err error) *Error {
	e.cause = err

	return e
}

func (e *Error) Causef(format string, args ...any) *Error {
	return e.Cause(fmt.Errorf(format, args...))
}

func (e *Error) Code(code int) *Error {
	if code < 0 || code > 9999 {
		panic(fmt.Sprintf("number out of range: %d", code))
	}

	e.code = fmt.Sprintf("E%04d", code)

	return e
}

func (e *Error) Help(help string) *Error {
	if help == "" {
		return e
	}

	e.helps = append(e.helps, help)

	return e
}

func (e *Error) Helpf(format string, args ...any) *Error {
	return e.Help(fmt.Sprintf(format, args...))
}

func (e *Error) HelpIf(help string, condition func() bool) *Error {
	if condition() {
		return e.Help(help)
	}

	return e
}

func (e *Error) HelpFunc(fn func() string) *Error {
	return e.Help(fn())
}

func (e *Error) Note(note string) *Error {
	if note == "" {
		return e
	}

	e.notes = append(e.notes, note)

	return e
}

func (e *Error) Notef(format string, args ...any) *Error {
	return e.Note(fmt.Sprintf(format, args...))
}

func (e *Error) AdditionalTemplateData(data map[string]any) *Error {
	e.additionalTemplateData = data

	return e
}

func (e *Error) GetMessage() string {
	return e.message
}

func (e *Error) Error() string {
	var result strings.Builder

	data := map[string]any{
		dataMessage: e.message,
		dataCause:   e.cause,
		dataWrapped: e.wrapped,
		dataCode:    e.code,
		dataNotes:   e.notes,
		dataHelps:   e.helps,
	}

	if len(e.additionalTemplateData) > 0 {
		maps.Copy(data, e.additionalTemplateData)
	}

	init := defaultInit
	if e.init != nil {
		init = e.init
	}

	if err := init.template.Execute(&result, data); err != nil {
		log.Printf("failed to execute error template: %v", err)
		// fall back to just the error message
		result.WriteString(e.message)
	}

	return result.String()
}

func (e *Error) WrappedErrors() []error {
	return e.wrapped
}
