package errors

import (
	"errors"
	"testing"

	"github.com/fatih/color"
	"github.com/stretchr/testify/assert"
)

func Test_Basic_Usage(t *testing.T) {
	original := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = original }()

	for _, tt := range []struct {
		name     string
		err      error
		expected string
	}{
		{
			name:     "no wrapped errors",
			err:      New("test"),
			expected: "error: test",
		},
		{
			name:     "no wrapped errors - formatted message",
			err:      Newf("formatted %s", "message"),
			expected: "error: formatted message",
		},
		{
			name:     "with error code",
			err:      New("test").Code(123),
			expected: "error[E0123]: test",
		},
		{
			name: "wrapped error: errors.Error",
			err: New("test").
				Wrap(New("wrapped")),
			expected: "error: test\n\nerror: wrapped",
		},
		{
			name: "cause",
			err: New("test").
				Cause(errors.New("cause of the 1st error\nplain error")).
				Wrap(New("wrapped").
					Cause(New("cause of the 2nd error\nerrors.New error"))),
			expected: `error: test
  --> cause of the 1st error
   | plain error

error: wrapped
  --> error: cause of the 2nd error
   | errors.New error`,
		},
		{
			name: "notes",
			err: New("notes").
				Cause(errors.New("cause of the 1st error\nplain error")).
				Note("this is because \nthis and this").
				Note("also \nthat"),
			expected: `error: notes
  --> cause of the 1st error
   | plain error
   = note: this is because 
           this and this
   = note: also 
           that`,
		},
		{
			name: "helps",
			err: New("helps").
				Cause(errors.New("cause of the 1st error\nplain error")).
				Help("do this \nbecause of reasons").
				Help("also \ndo that"),
			expected: `error: helps
  --> cause of the 1st error
   | plain error
   = help: do this 
           because of reasons
   = help: also 
           do that`,
		},
		{
			name: "help func",
			err: New("helps").
				Cause(errors.New("cause of the 1st error\nplain error")).
				HelpFunc(func() string { return "do this \nbecause of reasons" }).
				HelpFunc(func() string { return "" }), // this is empty so it does not appear in the output
			expected: `error: helps
  --> cause of the 1st error
   | plain error
   = help: do this 
           because of reasons`,
		},
		{
			name: "rust like error",
			err: New("'Foo' is not an iterator").
				Code(277).
				Causef(`src/main.rs:4:16

     for foo in Foo {}
                ^^^ 'Foo' is not an iterator
`).Note("maybe try calling '.iter()' or a similar method").
				Help("the trait 'std::iter::Iterator' is not implemented for 'Foo'").
				Note("required by 'std::iter::IntoIterator::into_iter'").
				Wrap(New("'&str' is not an iterator").
					Code(277).
					Causef(`src/main.rs:5:16

	 for foo in "" {}
				^^ '&str' is not an iterator
`).Help("call '.chars()' or '.bytes() on '&str'").
					Help("the trait 'std::iter::Iterator' is not implemented for '&str'").
					Note("required by 'std::iter::IntoIterator::into_iter'")),
			expected: `error[E0277]: 'Foo' is not an iterator
  --> src/main.rs:4:16
   | 
   |      for foo in Foo {}
   |                 ^^^ 'Foo' is not an iterator
   | 
   = note: maybe try calling '.iter()' or a similar method
   = note: required by 'std::iter::IntoIterator::into_iter'
   = help: the trait 'std::iter::Iterator' is not implemented for 'Foo'

error[E0277]: '&str' is not an iterator
  --> src/main.rs:5:16
   | 
   | 	 for foo in "" {}
   | 				^^ '&str' is not an iterator
   | 
   = note: required by 'std::iter::IntoIterator::into_iter'
   = help: call '.chars()' or '.bytes() on '&str'
   = help: the trait 'std::iter::Iterator' is not implemented for '&str'`,
		},
		{
			name: "extend",
			err: Extend(New("original").Causef("original cause")).
				Causef("override original cause"),
			expected: `error: original
  --> override original cause`,
		},
		{
			name: "extend with message",
			err: ExtendWithMessage(
				New("original").Causef("original cause"),
				"overridden message"),
			expected: `error: overridden message
  --> original cause`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func TestOutput_Global_Overrides(t *testing.T) {
	original := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = original }()

	for _, tt := range []struct {
		name     string
		init     func()
		err      error
		expected string
	}{
		{
			name: "override prefixMessage template",
			init: func() {
				SetMessagePrefixTemplate(`{{- define "messagePrefix" }}overridden{{ end }}`)
			},
			err:      New("test"),
			expected: `overridden: test`,
		},
		{
			name: "override cause template",
			init: func() {
				SetCauseTemplate(`{{- define "cause" }}
{{ print "    " .Cause }}{{ end }}`)
			},
			err: New("test").Cause(errors.New("xxx\n    yyy")),
			expected: `error: test
    xxx
    yyy`,
		},
		{
			name: "override cause template with a custom func",
			init: func() {
				AdditionalTemplateFunc("custom", func(input int) string {
					if input == 1 {
						return "one"
					}

					return "other"
				})
				SetCauseTemplate(`{{- define "cause" }}
1={{ custom 1 }}{{ end }}`)
			},
			err:      New("test").Cause(errors.New("xxx\n    yyy")),
			expected: "error: test\n1=one",
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			Reset()
			tt.init()

			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func Test_Reset(t *testing.T) {
	original := color.NoColor
	color.NoColor = true
	defer func() { color.NoColor = original }()

	originalFuncSize := len(funcMap)

	AdditionalTemplateFunc("custom", func() string { return "" })

	Reset()

	// the original funcMap stays the same
	assert.Equal(t, len(funcMap), originalFuncSize)
}
