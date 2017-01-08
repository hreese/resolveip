package resolveip

import (
	"bytes"
	"net"
	"strings"
)

type OutputConfig struct {
	Nonmatch          TextMutator
	ResolvedMatch     TextMutator
	UnresolvableMatch TextMutator
	Result            TextMutator
}

type ResolverFunc func(string, [][]int) string

func MakeResolver(out OutputConfig) ResolverFunc {
	return func(line string, matches [][]int) string {
		if len(matches) == 0 || len(line) == 0 {
			return line
		}
		var (
			buffer        bytes.Buffer
			nonmatchstart int = 0
		)

		for _, m := range matches {
			// print before first match
			buffer.WriteString(out.Nonmatch(line[nonmatchstart:m[0]]))
			nonmatchstart = m[1]

			match := line[m[0]:m[1]]
			// resolve and output match
			resolved, err := net.LookupAddr(match)
			if err == nil {
				buffer.WriteString(out.ResolvedMatch(match))
				buffer.WriteString(out.Result(strings.Join(resolved, ", ")))
			} else {
				buffer.WriteString(out.UnresolvableMatch(match))
			}
		}
		// print after last match
		buffer.WriteString(out.Nonmatch(line[matches[len(matches)-1][1]:len(line)]))

		return buffer.String()
	}
}
