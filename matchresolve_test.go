package resolveip

import (
	"bytes"
	"errors"
	"io"
	"testing"
)

func TestMakeResolverEmptyMatches(t *testing.T) {
	resolver := MakeResolver(OutputConfig{
		Nonmatch:          NOP,
		ResolvedMatch:     NOP,
		UnresolvableMatch: NOP,
		Result:            NOP,
	})

	line := "some line with text"
	got := resolver(line, nil)
	if got != line {
		t.Errorf("resolver with nil matches: got %q, want %q", got, line)
	}

	got = resolver(line, [][]int{})
	if got != line {
		t.Errorf("resolver with empty matches: got %q, want %q", got, line)
	}
}

func TestMakeResolverEmptyLine(t *testing.T) {
	resolver := MakeResolver(OutputConfig{
		Nonmatch:          NOP,
		ResolvedMatch:     NOP,
		UnresolvableMatch: NOP,
		Result:            NOP,
	})

	got := resolver("", [][]int{{0, 0}})
	if got != "" {
		t.Errorf("resolver with empty line: got %q, want %q", got, "")
	}
}

func TestMakeResolverUnresolvableFormatting(t *testing.T) {
	// Use an address that will not resolve (RFC 5737 TEST-NET).
	// net.LookupAddr will fail, so UnresolvableMatch is applied.
	tagger := func(tag string) TextMutator {
		return func(s string) string { return "[" + tag + ":" + s + "]" }
	}

	resolver := MakeResolver(OutputConfig{
		Nonmatch:          tagger("N"),
		ResolvedMatch:     tagger("R"),
		UnresolvableMatch: tagger("U"),
		Result:            tagger("RES"),
	})

	// Simulate a line "before 192.0.2.1 after" with a match on the IP part.
	line := "before 192.0.2.1 after"
	matches := [][]int{{7, 16}} // "192.0.2.1"

	got := resolver(line, matches)

	// The IP 192.0.2.1 (TEST-NET) should fail to resolve, so we expect:
	// [N:before ][U:192.0.2.1][N: after]
	want := "[N:before ][U:192.0.2.1][N: after]"
	if got != want {
		t.Errorf("resolver unresolvable:\n got %q\nwant %q", got, want)
	}
}

func TestMakeResolverNonmatchSegments(t *testing.T) {
	// Test with multiple matches to verify nonmatch segments between them.
	resolver := MakeResolver(OutputConfig{
		Nonmatch:          GenQuoter("{", "}"),
		ResolvedMatch:     NOP,
		UnresolvableMatch: GenQuoter("<", ">"),
		Result:            NOP,
	})

	line := "aaa 192.0.2.1 bbb 192.0.2.2 ccc"
	matches := [][]int{{4, 13}, {18, 27}} // two IPs

	got := resolver(line, matches)

	// Both IPs are unresolvable (TEST-NET), so:
	// {aaa }<192.0.2.1>{ bbb }<192.0.2.2>{ ccc}
	want := "{aaa }<192.0.2.1>{ bbb }<192.0.2.2>{ ccc}"
	if got != want {
		t.Errorf("resolver multiple matches:\n got %q\nwant %q", got, want)
	}
}

// --- InfoReader tests ---

func TestNewInfoReader(t *testing.T) {
	var buf bytes.Buffer
	msg := "hello info"
	ir := NewInfoReader(msg, &buf)

	p := make([]byte, 64)
	n, err := ir.Read(p)

	if n != 0 {
		t.Errorf("Read returned n=%d, want 0", n)
	}
	if !errors.Is(err, io.EOF) {
		t.Errorf("Read returned err=%v, want io.EOF", err)
	}
	if buf.String() != msg {
		t.Errorf("written message = %q, want %q", buf.String(), msg)
	}
}

func TestNewInfoReaderEmpty(t *testing.T) {
	var buf bytes.Buffer
	ir := NewInfoReader("", &buf)

	p := make([]byte, 64)
	n, err := ir.Read(p)

	if n != 0 {
		t.Errorf("Read returned n=%d, want 0", n)
	}
	if !errors.Is(err, io.EOF) {
		t.Errorf("Read returned err=%v, want io.EOF", err)
	}
	if buf.String() != "" {
		t.Errorf("written message = %q, want %q", buf.String(), "")
	}
}
