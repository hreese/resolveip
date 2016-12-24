package main

import (
	"bufio"
	"bytes"
	"fmt"
	"github.com/fatih/color"
	gar "github.com/goanywhere/regex"
	"net"
	"os"
	"regexp"
	"strings"
)

var (
	matchV4           = gar.IPv4.String()
	matchV6           = gar.IPv6.String()
	ip46              = regexp.MustCompile(fmt.Sprintf("(?:%s)|(?:%s)", matchV6, matchV4))
	hilightResolved   = color.New(color.FgGreen).SprintFunc()
	hilightUnresolved = color.New(color.FgRed).SprintFunc()
	hilightIP         = color.New(color.Bold).SprintFunc()
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
			buffer.WriteString(hilightIP(match))
			buffer.WriteString(hilightResolved(" »", strings.Join(resolved, ", "), "« "))
		} else {
			buffer.WriteString(hilightUnresolved(match))
		}
	}
	// print after last match
	buffer.WriteString(line[matches[len(matches)-1][1]:len(line)])

	return buffer.String()
}

func main() {
	// read input line by line
	scanner := bufio.NewScanner(os.Stdin)
	for scanner.Scan() {
		line := scanner.Text()
		// find all ip addresses
		matches := ip46.FindAllStringIndex(scanner.Text(), -1)
		fmt.Println(resolveIPs(line, matches))
	}
	if err := scanner.Err(); err != nil {
		fmt.Fprintln(os.Stderr, "reading standard input:", err)
	}
}
