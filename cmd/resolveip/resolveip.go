//go:generate goversioninfo -icon=../../res/icon.ico

// Install goversioninfo to generate and embed a Windows icon:
//
//	go install github.com/josephspurrier/goversioninfo/cmd/goversioninfo@latest

package main

import (
	"bufio"
	"flag"
	"fmt"
	"io"
	"os"
	"runtime"

	"github.com/fatih/color"
	"github.com/hreese/goodregex"
	resolveip "github.com/hreese/resolveip/v2"
)

// openInputs builds an [io.Reader] from the file arguments on the command line.
// If no files are given, it returns [os.Stdin]. On Windows (without -batch),
// stdin is appended after the files so the console stays open.
func openInputs(batch bool) io.Reader {
	if flag.NArg() == 0 {
		return os.Stdin
	}

	infiles := make([]io.Reader, 0, flag.NArg())
	for _, filename := range flag.Args() {
		reader, err := os.Open(filename)
		if err != nil {
			fmt.Fprintf(os.Stderr, "error opening %s: %v\n", filename, err)
			continue
		}
		infiles = append(infiles, reader)
	}

	// On Windows, keep stdin open after processing files unless -batch is set.
	// This prevents the console window from closing after drag-and-drop.
	if runtime.GOOS == "windows" && !batch {
		return io.MultiReader(io.MultiReader(infiles...),
			resolveip.NewInfoReader("# Start typing or pasting text here\n", os.Stderr),
			os.Stdin)
	}
	return io.MultiReader(infiles...)
}

func main() {
	// Parse command-line flags.
	confWantColor := flag.Bool("c", false, "Enforce ANSI color codes")
	confNoColor := flag.Bool("C", false, "Disable ANSI color codes")
	no4 := flag.Bool("no4", false, "Disable matching of IPv4 addresses")
	no6 := flag.Bool("no6", false, "Disable matching of IPv6 addresses")
	confBatch := flag.Bool("batch", false, "Don't read from stdin after all files are processed (Windows only)")

	flag.Parse()

	// Apply color settings after flag.Parse().
	confMatchV4 := !*no4
	confMatchV6 := !*no6

	if *confNoColor {
		color.NoColor = true
	}
	if *confWantColor {
		color.NoColor = false
	}

	// Configure output formatting.
	outputConfig := resolveip.OutputConfig{
		Nonmatch:          resolveip.GenHighlighter(color.Faint),
		UnresolvableMatch: resolveip.Chain(resolveip.GenHighlighter(color.FgRed)),
		ResolvedMatch:     resolveip.Chain(resolveip.GenHighlighter(color.Bold)),
		Result:            resolveip.Chain(resolveip.GenQuoter(" \u00bb", "\u00ab "), resolveip.GenHighlighter(color.FgGreen)),
	}
	resolver := resolveip.MakeResolver(outputConfig)

	// Process input line by line.
	input := openInputs(*confBatch)
	scanner := bufio.NewScanner(input)
	for scanner.Scan() {
		line := scanner.Text()
		if confMatchV6 {
			matches := goodregex.MatchIPv6.FindAllStringIndex(line, -1)
			line = resolver(line, matches)
		}
		if confMatchV4 {
			matches := goodregex.MatchIPv4.FindAllStringIndex(line, -1)
			line = resolver(line, matches)
		}
		fmt.Fprintln(color.Output, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading input:", err)
	}
}
