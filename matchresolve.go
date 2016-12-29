package resolveip

import (
	"bytes"
	"net"
	"regexp"
	"strings"
)

const (
	regexMatchV6 = `(
        # addresses containing ::
                                                      :: (?:(?:[[:alnum:]]{1,4}:){0,6}(?:[[:alnum:]]{1,4})){0,1} | 
                                 (?:[[:alnum:]]{1,4}) :: (?:(?:[[:alnum:]]{1,4}:){0,5}(?:[[:alnum:]]{1,4})){0,1} |
        (?:[[:alnum:]]{1,4}:){1} (?:[[:alnum:]]{1,4}) :: (?:(?:[[:alnum:]]{1,4}:){0,4}(?:[[:alnum:]]{1,4})){0,1} |
        (?:[[:alnum:]]{1,4}:){2} (?:[[:alnum:]]{1,4}) :: (?:(?:[[:alnum:]]{1,4}:){0,3}(?:[[:alnum:]]{1,4})){0,1} |
        (?:[[:alnum:]]{1,4}:){3} (?:[[:alnum:]]{1,4}) :: (?:(?:[[:alnum:]]{1,4}:){0,2}(?:[[:alnum:]]{1,4})){0,1} |
        (?:[[:alnum:]]{1,4}:){4} (?:[[:alnum:]]{1,4}) :: (?:(?:[[:alnum:]]{1,4}:){0,1}(?:[[:alnum:]]{1,4})){0,1} |
        (?:[[:alnum:]]{1,4}:){5} (?:[[:alnum:]]{1,4}) :: (?:[[:alnum:]]{1,4}){0,1}                               |
        (?:[[:alnum:]]{1,4}:){6} (?:[[:alnum:]]{1,4}) ::                                                         |
        # plain IPv6 address
        (?:[[:alnum:]]{1,4}:){7} (?:[[:alnum:]]{1,4})
        )`
	regexMatchV4 = `
        (?:(?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.){3}
        (?:25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`
)

type OutputConfig struct {
	Nonmatch          TextMutator
	ResolvedMatch     TextMutator
	UnresolvableMatch TextMutator
	Result            TextMutator
}

var (
	regexRemoveComments = regexp.MustCompile("(?m:#.*$)")
	regexRemoveSpaces   = regexp.MustCompile("(?m:[[:space:]]+)")
	MatchV4             = MustCompileReadableRegex(regexMatchV4)
	MatchV6             = MustCompileReadableRegex(regexMatchV6)
)

// MustCompileReadableRegex fixes and compiles a readable regex
func MustCompileReadableRegex(r string) *regexp.Regexp {
	r = regexRemoveComments.ReplaceAllString(r, "")
	r = regexRemoveSpaces.ReplaceAllString(r, "")
	return regexp.MustCompile(r)
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
