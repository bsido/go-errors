package warnings

import (
	"testing"

	"github.com/bsido/go-errors/errors"
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
			name:     "warning",
			err:      New("test"),
			expected: "warning: test",
		},
		{
			name:     "warning format",
			err:      Newf("test %s", "format"),
			expected: "warning: test format",
		},
		{
			name: "warning with error cause",
			err:  New("wrapper").Cause(errors.New("error")),
			expected: `warning: wrapper
  --> error: error`,
		},
		{
			name: "warning as a wrapper",
			err:  New("wrapper").Wrap(errors.New("error")),
			expected: `warning: wrapper

error: error`,
		},
	} {
		t.Run(tt.name, func(t *testing.T) {
			assert.Equal(t, tt.expected, tt.err.Error())
		})
	}
}

func Test_Is(t *testing.T) {
	warning := New("test")
	assert.True(t, Is(warning))

	warning = New("wrapper").Cause(errors.New("error"))
	assert.True(t, Is(warning))
}
