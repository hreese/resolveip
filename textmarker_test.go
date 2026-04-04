package resolveip

import (
	"testing"

	"github.com/fatih/color"
)

func TestNOP(t *testing.T) {
	cases := []string{"", "hello", "hello world", "123.45.67.89"}
	for _, input := range cases {
		if got := NOP(input); got != input {
			t.Errorf("NOP(%q) = %q, want %q", input, got, input)
		}
	}
}

func TestDelete(t *testing.T) {
	cases := []string{"", "hello", "some longer text"}
	for _, input := range cases {
		if got := Delete(input); got != "" {
			t.Errorf("Delete(%q) = %q, want %q", input, got, "")
		}
	}
}

func TestGenQuoter(t *testing.T) {
	tests := []struct {
		before, after string
		input         string
		want          string
	}{
		{"[", "]", "hello", "[hello]"},
		{"(", ")", "world", "(world)"},
		{"<", ">", "", "<>"},
		{"", "", "text", "text"},
		{"<<", ">>", "ip", "<<ip>>"},
	}
	for _, tc := range tests {
		q := GenQuoter(tc.before, tc.after)
		if got := q(tc.input); got != tc.want {
			t.Errorf("GenQuoter(%q, %q)(%q) = %q, want %q",
				tc.before, tc.after, tc.input, got, tc.want)
		}
	}
}

func TestGenHighlighter(t *testing.T) {
	h := GenHighlighter(color.FgRed)
	got := h("test")
	if got == "" {
		t.Error("GenHighlighter(FgRed)(\"test\") returned empty string")
	}

	// With color disabled, SprintFunc returns the original text.
	// With color enabled, the result contains ANSI codes wrapping "test".
	// Either way the original text must be present in the output.
	if len(got) < len("test") {
		t.Errorf("GenHighlighter output %q is shorter than input %q", got, "test")
	}
}

func TestChainEmpty(t *testing.T) {
	c := Chain()
	if got := c("hello"); got != "hello" {
		t.Errorf("Chain()(%q) = %q, want %q", "hello", got, "hello")
	}
}

func TestChainSingle(t *testing.T) {
	c := Chain(Delete)
	if got := c("hello"); got != "" {
		t.Errorf("Chain(Delete)(%q) = %q, want %q", "hello", got, "")
	}
}

func TestChainMultiple(t *testing.T) {
	// Apply two quoters in sequence: first wraps with [], then with ()
	q1 := GenQuoter("[", "]")
	q2 := GenQuoter("(", ")")
	c := Chain(q1, q2)
	want := "([hello])"
	if got := c("hello"); got != want {
		t.Errorf("Chain(q1, q2)(%q) = %q, want %q", "hello", got, want)
	}
}

func TestChainOrder(t *testing.T) {
	// Verify that chain applies mutators left-to-right
	upper := func(_ string) string { return "UPPER" }
	wrap := GenQuoter("<", ">")

	// upper first, then wrap → "<UPPER>"
	c1 := Chain(upper, wrap)
	if got := c1("x"); got != "<UPPER>" {
		t.Errorf("Chain(upper, wrap)(\"x\") = %q, want %q", got, "<UPPER>")
	}

	// wrap first, then upper → "UPPER" (upper ignores input)
	c2 := Chain(wrap, upper)
	if got := c2("x"); got != "UPPER" {
		t.Errorf("Chain(wrap, upper)(\"x\") = %q, want %q", got, "UPPER")
	}
}
