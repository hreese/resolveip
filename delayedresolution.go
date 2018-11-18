package resolveip

import (
	"bytes"
	"net"
	"strings"
)

type FutureResult func() string
type ResolutionResult struct {
	names []string
	err   error
}

func AsyncResolve(out OutputConfig, line string, matches [][]int) []FutureResult {
	var results []FutureResult

	if len(matches) == 0 || len(line) == 0 {
		return []FutureResult{func() string { return line }}
	}
	var (
		nonmatchstart int = 0
	)
	for _, m := range matches {
		// print before first match
		results = append(results, func() string {
			return out.Nonmatch(line[nonmatchstart:m[0]])
		})
		nonmatchstart = m[1]

		match := line[m[0]:m[1]]
		// resolve
		var resolutionchan = make(chan ResolutionResult)
		go func(c chan ResolutionResult) {
			var r ResolutionResult
			r.names, r.err = net.LookupAddr(match)
			c <- r
			close(c)
		}(resolutionchan)
		results = append(results, func() string {
			var buffer bytes.Buffer
			resolved := <-resolutionchan
			if resolved.err == nil {
				buffer.WriteString(out.ResolvedMatch(match))
				buffer.WriteString(out.Result(strings.Join(resolved.names, ", ")))
			} else {
				buffer.WriteString(out.UnresolvableMatch(match))
			}
			return buffer.String()
		})
	}
	// print after last match
	results = append(results, func() string {
		return out.Nonmatch(line[matches[len(matches)-1][1]:])
	})
	return results
}
