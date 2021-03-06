//go:generate goversioninfo -icon=../../res/icon.ico

// install goversioninfo to generate and embed a Windows Icon
//   go get github.com/josephspurrier/goversioninfo/cmd/goversioninfo

package main

import (
	"bufio"
	"flag"
	"fmt"
	"github.com/fatih/color"
	. "github.com/hreese/goodregex"
	. "github.com/hreese/resolveip"
	"io"
	"os"
	"runtime"
)

var (
	// default output config
	outputConfig = OutputConfig{
		Nonmatch:          GenHighlighter(color.Faint),
		UnresolvableMatch: Chain(GenHighlighter(color.FgRed)),
		ResolvedMatch:     Chain(GenHighlighter(color.Bold)),
		Result:            Chain(GenQuoter(" »", "« "), GenHighlighter(color.FgGreen)),
	}
	confMatchV4   bool
	confMatchV6   bool
	confBatch     bool
	confWantColor *bool
	confNoColor   *bool
	resolveIPs    ResolverFunc
)

// parse commandline arguments
func init() {
	confWantColor = flag.Bool("c", false, "Enforce ANSI color codes")
	confNoColor = flag.Bool("C", false, "Disable ANSI color codes")
	confMatchV4 = !*flag.Bool("no4", false, "Disable matching of IPv4 addresses")
	confMatchV6 = !*flag.Bool("no6", false, "Disable matching of IPv6 addresses")
	flag.BoolVar(&confBatch, "batch", false, "Does not read from stdin after all files are processed (Windows only)")

	flag.Parse()

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
		// Use case for this is: user drops text file(s) onto executable
		// on Windows, this prevents the program windows to close upon exit
		if runtime.GOOS == "windows" && confBatch == false {
			input = io.MultiReader(io.MultiReader(infiles...),
				NewInfoWriter("# Start typing or pasting text here\n", os.Stderr),
				os.Stdin)
		} else {
			input = io.MultiReader(infiles...)
		}
	}
	// read input line by line
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		// find all IPv6 addresses
		if confMatchV6 {
			matches := MatchIPv6.FindAllStringIndex(line, -1)
			line = resolveIPs(line, matches)
		}
		// find all IPv4 addresses
		if confMatchV4 {
			matches := MatchIPv4.FindAllStringIndex(line, -1)
			line = resolveIPs(line, matches)
		}
		fmt.Fprintln(color.Output, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}

}
