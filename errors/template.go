package errors

import (
	"strings"
	"text/template"

	"github.com/fatih/color"
)

const (
	templateNameError = "error"

	funcBold      = "bold"
	funcBoldRed   = "boldRed"
	funcBoldBlue  = "boldBlue"
	funcBoldGreen = "boldGreen"
	funcSplit     = "split"

	dataMessage = "Message"
	dataWrapped = "Wrapped"
	dataCode    = "Code"
	dataCause   = "Cause"
	dataNotes   = "Notes"
	dataHelps   = "Helps"
)

type TemplateDefinition string

const (
	TemplateDefinitionCause         TemplateDefinition = "cause"
	TemplateDefinitionMessagePrefix TemplateDefinition = "messagePrefix"
	TemplateDefinitionNotes         TemplateDefinition = "notes"
	TemplateDefinitionHelps         TemplateDefinition = "helps"
)

const (
	errorTemplate = `{{- template "messagePrefix" . }}{{- bold (print ": " .Message) }}
{{- template "cause" . }}

{{- template "notes" . }}

{{- template "helps" . }}

{{- if .Wrapped }}
{{- range .Wrapped }}

{{ . }}
{{- end }}
{{- end }}`

	causeTemplate = `{{- define "cause" }}
{{- if .Cause }}
{{- $cause := split (print .Cause) "\n" }}
  {{ boldBlue "--> " -}}{{ index $cause 0 }}
  {{- range $line := slice $cause 1 }}
   {{ boldBlue "| " }}{{ . }}
  {{- end }}
{{- end }}
{{- end }}`

	messagePrefixTemplate = `{{- define "messagePrefix" }}
	{{- if .Code }}
		{{- boldRed (print "error[" .Code "]") }}
	{{- else }}
		{{- boldRed "error" }}
	{{- end }}
{{- end }}`

	notesTemplate = `{{- define "notes" }}
{{- if .Notes }}
   {{- range $note := .Notes }}
   {{- $lines := split $note "\n" }}
   {{ boldBlue "= " }}{{ bold "note" }}: {{ index $lines 0 -}}
       {{- range slice $lines 1 }}
           {{ . }}
       {{- end }}
   {{- end }}
{{- end }}
{{- end }}`

	helpsTemplate = `{{- define "helps" }}
{{- if .Helps }}
   {{- range $help := .Helps }}
   {{- $lines := split $help "\n" }}
   {{ boldBlue "= " }}{{ boldGreen "help" }}: {{ index $lines 0 -}}
       {{- range slice $lines 1 }}
           {{ . }}
       {{- end }}
   {{- end }}
{{- end }}
{{- end }}`
)

var (
	bold      = color.New(color.Bold).Sprintf
	boldRed   = color.New(color.FgRed, color.Bold).Sprintf
	boldBlue  = color.New(color.FgBlue, color.Bold).Sprintf
	boldGreen = color.New(color.FgGreen, color.Bold).Sprintf
)

var funcMap = template.FuncMap{
	funcBold:      bold,
	funcBoldRed:   boldRed,
	funcBoldGreen: boldGreen,
	funcBoldBlue:  boldBlue,
	funcSplit:     strings.Split,
}

var defaultInit *Init

func init() {
	Reset()
}

func Reset() {
	defaultInit = NewInitializer()
}

func SetCauseTemplate(template string) {
	defaultInit.definitions[TemplateDefinitionCause] = template
	defaultInit.template = newTemplate(defaultInit.definitions, defaultInit.funcMap)
}

func SetMessagePrefixTemplate(template string) {
	defaultInit.definitions[TemplateDefinitionMessagePrefix] = template
	defaultInit.template = newTemplate(defaultInit.definitions, defaultInit.funcMap)
}

func SetNotesTemplate(template string) {
	defaultInit.definitions[TemplateDefinitionNotes] = template
	defaultInit.template = newTemplate(defaultInit.definitions, defaultInit.funcMap)
}

func SetHelpsTemplate(template string) {
	defaultInit.definitions[TemplateDefinitionHelps] = template
	defaultInit.template = newTemplate(defaultInit.definitions, defaultInit.funcMap)
}

func AdditionalTemplateFunc(name string, fn any) {
	AdditionalTemplateFuncs(template.FuncMap{name: fn})
}

func AdditionalTemplateFuncs(funcs template.FuncMap) {
	functions := defaultInit.funcMap

	for name, fn := range funcs {
		functions[name] = fn
	}

	defaultInit.funcMap = functions
	defaultInit.template = newTemplate(defaultInit.definitions, functions)
}

func newTemplate(definitions map[TemplateDefinition]string, fns template.FuncMap) *template.Template {
	result := template.New(templateNameError).Funcs(fns)

	for _, def := range definitions {
		result = template.Must(result.Parse(def))
	}

	return template.Must(result.Parse(errorTemplate))
}
