package resolveip

import (
	"github.com/fatih/color"
)

// TextMutator changes a given string
type TextMutator func(string) string

// Chain returns a new TextMutator that applies all given TextMutator functions in order
func Chain(mfuncs ...TextMutator) TextMutator {
	return func(input string) string {
		for _, m := range mfuncs {
			input = m(input)
		}
		return input
	}
}

// Delete returns an empty string
func Delete(input string) string {
	return ""
}

// GenQuoter returns a new TextMutator that surrounds the input string with custom strings
func GenQuoter(before, after string) TextMutator {
	return func(input string) string {
		return before + input + after
	}
}

// GenHighlighter returns a new TextMutator that colorizes the input string
func GenHighlighter(attr ...color.Attribute) TextMutator {
	c := color.New(attr...)
	sprint := c.SprintFunc()

	return func(input string) string {
		return sprint(input)
	}
}

// NOP simply returns the input string
func NOP(input string) string {
	return input
}
