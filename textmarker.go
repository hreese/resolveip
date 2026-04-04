package resolveip

import (
	"github.com/fatih/color"
)

// TextMutator is a function that transforms a string. It is used to apply
// formatting (colors, quoting, etc.) to text segments in resolver output.
type TextMutator func(string) string

// Chain returns a new [TextMutator] that applies all given mutators in order.
func Chain(mfuncs ...TextMutator) TextMutator {
	return func(input string) string {
		for _, m := range mfuncs {
			input = m(input)
		}
		return input
	}
}

// Delete returns an empty string, effectively removing the input text.
func Delete(_ string) string {
	return ""
}

// GenQuoter returns a [TextMutator] that surrounds input text with the
// specified before and after strings.
func GenQuoter(before, after string) TextMutator {
	return func(input string) string {
		return before + input + after
	}
}

// GenHighlighter returns a [TextMutator] that colorizes input text using
// the specified [color.Attribute] values.
func GenHighlighter(attr ...color.Attribute) TextMutator {
	c := color.New(attr...)
	sprint := c.SprintFunc()

	return func(input string) string {
		return sprint(input)
	}
}

// NOP returns the input string unchanged.
func NOP(input string) string {
	return input
}
