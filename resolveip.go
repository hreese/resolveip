package main

import (
	"bytes"
	"bufio"
	"fmt"
	"github.com/fatih/color"
	"net"
	"os"
	"regexp"
    "strings"
)

var (
	// source: https://github.com/goanywhere/regex/blob/master/regex.go
    // FIXME: allow for ›::‹ in IPv6-Addr
	ip46              = regexp.MustCompile(`(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)\.(25[0-5]|2[0-4][0-9]|[01]?[0-9][0-9]?)|(([0-9a-fA-F]{1,4}:){7,7}[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,7}:|([0-9a-fA-F]{1,4}:){1,6}:[0-9a-fA-F]{1,4}|([0-9a-fA-F]{1,4}:){1,5}(:[0-9a-fA-F]{1,4}){1,2}|([0-9a-fA-F]{1,4}:){1,4}(:[0-9a-fA-F]{1,4}){1,3}|([0-9a-fA-F]{1,4}:){1,3}(:[0-9a-fA-F]{1,4}){1,4}|([0-9a-fA-F]{1,4}:){1,2}(:[0-9a-fA-F]{1,4}){1,5}|[0-9a-fA-F]{1,4}:((:[0-9a-fA-F]{1,4}){1,6})|:((:[0-9a-fA-F]{1,4}){1,7}|:)|fe80:(:[0-9a-fA-F]{0,4}){0,4}%[0-9a-zA-Z]{1,}|::(ffff(:0{1,4}){0,1}:){0,1}((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])|([0-9a-fA-F]{1,4}:){1,4}:((25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9])\.){3,3}(25[0-5]|(2[0-4]|1{0,1}[0-9]){0,1}[0-9]))`)
	hilightResolved   = color.New(color.FgGreen).SprintFunc()
	hilightUnresolved = color.New(color.FgRed).SprintFunc()
	hilightIP         = color.New(color.Bold).SprintFunc()
)

func resolveIPs(line string, matches [][]int) string {
	if len(matches) == 0 || len(line) == 0 {
		return line
	}
	var (
        buffer bytes.Buffer
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
