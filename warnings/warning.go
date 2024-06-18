package warnings

import (
	goerrors "errors"
	"fmt"
	"strings"
	"text/template"

	"github.com/bsido/go-errors/errors"
	"github.com/fatih/color"
)

const (
	funcBoldYellow = "boldYellow"

	messagePrefixTemplate = `{{- define "messagePrefix" }}
	{{- if .Code }}
		{{- boldYellow (print "warning[" .Code "]") }}
	{{- else }}
		{{- boldYellow "warning" }}
	{{- end }}
{{- end }}`
)

var (
	boldYellow = color.New(color.FgYellow, color.Bold).Sprintf
)

var warningsInit *errors.Init

func init() {
	warningsInit = errors.NewInitializer(
		errors.WithTemplateDefinition(errors.TemplateDefinitionMessagePrefix, messagePrefixTemplate),
		errors.WithAdditionalFunctions(template.FuncMap{
			funcBoldYellow: boldYellow,
		}))
}

func New(message string) *errors.Error {
	return warningsInit.NewError(message)
}

func Newf(format string, args ...any) *errors.Error {
	return New(fmt.Sprintf(format, args...))
}

func From(err error) *errors.Error {
	var e *errors.Error
	if goerrors.As(err, &e) {
		return warningsInit.Extend(e)
	}

	return New(err.Error())
}

// Is returns true if the error is a warning
func Is(original error) bool {
	var err *errors.Error

	ok := goerrors.As(original, &err)
	if !ok {
		return false
	}

	if strings.Contains(strings.SplitN(err.Error(), "\n", 2)[0], "warning") {
		return true
	}

	return false
}
