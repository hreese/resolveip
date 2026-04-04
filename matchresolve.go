package resolveip

import (
	"bytes"
	"net"
	"strings"
)

// OutputConfig controls how matched and unmatched text segments are formatted
// in the resolver output.
type OutputConfig struct {
	Nonmatch          TextMutator // Applied to text between matches.
	ResolvedMatch     TextMutator // Applied to IPs that resolve successfully.
	UnresolvableMatch TextMutator // Applied to IPs that fail to resolve.
	Result            TextMutator // Applied to the resolved hostname(s).
}

// ResolverFunc processes a line of text and a set of regex match indices,
// resolving any matched IP addresses and returning the formatted result.
type ResolverFunc func(string, [][]int) string

// MakeResolver returns a [ResolverFunc] that resolves IP addresses found at
// the given match positions using [net.LookupAddr] and formats the output
// according to the provided [OutputConfig].
func MakeResolver(out OutputConfig) ResolverFunc {
	return func(line string, matches [][]int) string {
		if len(matches) == 0 || len(line) == 0 {
			return line
		}

		var (
			buffer        bytes.Buffer
			nonmatchStart int
		)

		for _, m := range matches {
			// Write text before this match.
			buffer.WriteString(out.Nonmatch(line[nonmatchStart:m[0]]))
			nonmatchStart = m[1]

			match := line[m[0]:m[1]]
			// Resolve and format the matched IP address.
			resolved, err := net.LookupAddr(match)
			if err == nil {
				buffer.WriteString(out.ResolvedMatch(match))
				buffer.WriteString(out.Result(strings.Join(resolved, ", ")))
			} else {
				buffer.WriteString(out.UnresolvableMatch(match))
			}
		}
		// Write text after the last match.
		buffer.WriteString(out.Nonmatch(line[nonmatchStart:]))

		return buffer.String()
	}
}
