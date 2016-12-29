package main

import (
	"bufio"
	"flag"
	"fmt"
	. "git.heiko-reese.de/hreese/resolveip"
    "github.com/fatih/color"
	"io"
	"os"
)

var (
	// default output config
	outputConfig = OutputConfig{
		Nonmatch:          GenHighlighter(color.Faint),
		UnresolvableMatch: Chain(GenHighlighter(color.FgRed)),
		ResolvedMatch:     Chain(GenHighlighter(color.Bold)),
		Result:            Chain(GenQuoter(" »", "« "), GenHighlighter(color.FgGreen)),
	}
        confMatchV4 bool
        confMatchV6 bool
	resolveIPs ResolverFunc
)

// parse commandline arguments
func init() {
	var (
		confWantColor = flag.Bool("c", false, "Enforce ANSI color codes")
		confNoColor   = flag.Bool("C", false, "Disable ANSI color codes")
		no4           = flag.Bool("no4", false, "Disable ANSI color codes")
		no6           = flag.Bool("no6", false, "Disable ANSI color codes")
	)

	flag.Parse()

	confMatchV4 = !*no4
	confMatchV6 = !*no6
	if *confNoColor {
		color.NoColor = true
	}
	if *confWantColor {
		color.NoColor = false
	}

	resolveIPs = MakeResolver(outputConfig)
}

func main() {
	var input io.Reader

	if flag.NArg() == 0 {
		// read from stdin if no files are given
		input = os.Stdin
	} else {
		// read files from argument list
		var infiles []io.Reader
		for _, filename := range flag.Args() {
			reader, err := os.Open(filename)
			if err == nil {
				infiles = append(infiles, reader)
			}
		}
		input = io.MultiReader(infiles...)
	}
	// read input line by line
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		// find all IPv6 addresses
		if confMatchV6 {
			matches := MatchV6.FindAllStringIndex(line, -1)
			line = resolveIPs(line, matches)
		}
		// find all IPv4 addresses
		if confMatchV4 {
			matches := MatchV4.FindAllStringIndex(line, -1)
			line = resolveIPs(line, matches)
		}
		fmt.Fprintln(color.Output, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
