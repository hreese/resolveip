package main

import (
	"bufio"
	"bytes"
	"flag"
	"fmt"
	"github.com/fatih/color"
	"io"
	"net"
	"os"
	"regexp"
	"strings"
)

var (
	matchV4 = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)`)
	// TODO: this still fails to match xxx::1
	matchV6             = regexp.MustCompile(`(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`)
	highlightResolved   = color.New(color.FgGreen).SprintFunc()
	highlightUnresolved = color.New(color.FgRed).SprintFunc()
	highlightIP         = color.New(color.Bold).SprintFunc()
	confWantColor       bool
	confNoColor         bool
	confMatchV4         bool
	confMatchV6         bool
)

func resolveIPs(line string, matches [][]int) string {
	if len(matches) == 0 || len(line) == 0 {
		return line
	}
	var (
		buffer        bytes.Buffer
		nonmatchstart int = 0
	)

	for _, m := range matches {
		// print before first match
		buffer.WriteString(line[nonmatchstart:m[0]])
		nonmatchstart = m[1]

		match := line[m[0]:m[1]]
		// resolve and output match
		resolved, err := net.LookupAddr(match)
		if err == nil {
			buffer.WriteString(highlightIP(match))
			buffer.WriteString(highlightResolved(" »", strings.Join(resolved, ", "), "« "))
		} else {
			buffer.WriteString(highlightUnresolved(match))
		}
	}
	// print after last match
	buffer.WriteString(line[matches[len(matches)-1][1]:len(line)])

	return buffer.String()
}

func init() {
	flag.BoolVar(&confWantColor, "c", false, "Enforce ANSI color codes")
	flag.BoolVar(&confNoColor, "C", false, "Disable ANSI color codes")
	no4 := flag.Bool("no4", false, "Disable ANSI color codes")
	no6 := flag.Bool("no6", false, "Disable ANSI color codes")
	flag.Parse()
	confMatchV4 = !*no4
	confMatchV6 = !*no6
	if confNoColor {
		color.NoColor = true
	}
	if confWantColor {
		color.NoColor = false
	}
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
			matches := matchV6.FindAllStringIndex(line, -1)
			line = resolveIPs(line, matches)
		}
		// find all IPv4 addresses
		if confMatchV4 {
			matches := matchV4.FindAllStringIndex(line, -1)
			line = resolveIPs(line, matches)
		}
		fmt.Fprintln(color.Output, line)
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
