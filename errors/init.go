package errors

import (
	"fmt"
	"maps"
	"text/template"
)

type Init struct {
	funcMap  template.FuncMap
	template *template.Template

	definitions map[TemplateDefinition]string
}

func NewInitializer(opts ...InitOption) *Init {
	options := newOptions(opts)

	return &Init{
		funcMap:  options.funcMap,
		template: newTemplate(options.definitions, options.funcMap),

		definitions: options.definitions,
	}
}

func (b *Init) NewError(message string) *Error {
	return New(message).initializer(b)
}

func (b *Init) NewErrorf(format string, args ...any) *Error {
	return b.NewError(fmt.Sprintf(format, args...))
}

func (b *Init) Extend(original *Error) *Error {
	return original.initializer(b)
}

type templateOptions struct {
	funcMap     template.FuncMap
	definitions map[TemplateDefinition]string
}

func newOptions(opts []InitOption) *templateOptions {
	result := &templateOptions{
		funcMap: maps.Clone(funcMap),
		definitions: map[TemplateDefinition]string{
			TemplateDefinitionMessagePrefix: messagePrefixTemplate,
			TemplateDefinitionCause:         causeTemplate,
			TemplateDefinitionNotes:         notesTemplate,
			TemplateDefinitionHelps:         helpsTemplate,
		},
	}

	for _, opt := range opts {
		opt(result)
	}

	return result
}

type InitOption func(*templateOptions)

func WithFunctions(funcMap template.FuncMap) InitOption {
	return func(opts *templateOptions) {
		opts.funcMap = funcMap
	}
}

func WithAdditionalFunction(name string, fn any) InitOption {
	return func(opts *templateOptions) {
		opts.funcMap[name] = fn
	}
}

func WithAdditionalFunctions(funcs template.FuncMap) InitOption {
	return func(opts *templateOptions) {
		for name, fn := range funcs {
			opts.funcMap[name] = fn
		}
	}
}

func WithTemplateDefinition(name TemplateDefinition, definition string) InitOption {
	return func(opts *templateOptions) {
		opts.definitions[name] = definition
	}
}
