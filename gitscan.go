package main

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"os/exec"
	"strconv"
	"strings"
	"time"
)

type commit struct {
	alignment
	user                     *author
	hasEndingPeriod, isMerge bool
	date                     time.Time
	message                  string
}
type author struct {
	commits int
	alignment
	accumulator alignment
	name, email string
}

type gitOption func([]string) []string

var doNothing gitOption = func(s []string) []string { return s }

func optionNoMerges(none bool) gitOption {
	if none {
		return func(s []string) []string { return append(s, "--no-merges") }
	}
	return doNothing
}

func optionAuthorPattern(author string) gitOption {
	if author != "" {
		return func(s []string) []string { return append(s, "--author", author) }
	}
	return doNothing
}

func optionMaxCommits(n int) gitOption {
	return func(s []string) []string { return append(s, "-n", strconv.Itoa(n)) }
}

func optionBranch(b string) gitOption {
	if b == "" {
		b = "--all"
	}
	return func(s []string) []string { return append([]string{b}, s...) }
}

// Stats return author alignment in human readable format (with newlines)
func (a author) Stats() string {
	return fmt.Sprintf("Author %v is %v\nCommits: %d\nAccumulated:%0.1g\n",
		a.name, a.alignment.Format(), a.commits, a.accumulator)
}

// ScanCWD Scans .git in current working directory using git
// command. Scans up to maxCommit messages.
func ScanCWD(opts ...gitOption) ([]commit, []author, error) {
	var args []string
	for i := range opts {
		args = opts[i](args)
	}
	args = append([]string{"log"}, args...)
	cmd := exec.Command("git", args...)
	reader, writer := io.Pipe()
	cmd.Stdout = writer
	cmdstderr := &strings.Builder{}
	cmd.Stderr = cmdstderr
	go func() {
		cmd.Run()
		writer.Close()
	}()
	commits, authors, err := GitLogScan(reader)
	if err == io.EOF {
		err = nil
	}
	errmsg := cmdstderr.String()
	if err == nil && errmsg != "" {
		err = errors.New(errmsg)
	}
	return commits, authors, err
}

// GitLogScan reads git log results and generates commits
func GitLogScan(r io.Reader) (commits []commit, authors []author, err error) {
	rdr := bufio.NewReader(r)
	commits = make([]commit, 0, maxCommits)
	authors = make([]author, maxAuthors)
	authmap := make(map[string]*author)
	var c commit
	var a author
	var auth *author
	counter := 0
	eof := false
	for !eof {
		if counter >= maxCommits {
			break
		}
		c, a, err = scanNextCommit(rdr)
		if err == errSkipCommit {
			continue
		}
		if err == io.EOF {
			eof, err = true, nil
		}
		if err != nil {
			break
		}
		auth, err = processAuthor(a, authors, authmap)
		if err == errSkipCommit {
			continue
		}
		processCommit(&c, auth)
		commits = append(commits, c)
		counter++
	}
	if err == errSkipCommit || err == io.EOF {
		err = nil
	}
	return commits, authors[0:len(authmap)], err
}

func processAuthor(a author, authors []author, authmap map[string]*author) (*author, error) {
	// if author name is blank, then skip the person
	if a.name == "" {
		return nil, errSkipCommit
	}
	// find author in list
	author, ok := authmap[a.name]
	nAuthors := len(authmap)
	if !ok {
		if len(authmap) == len(authors) {
			return author, errSkipCommit
		}
		authors[nAuthors] = a
		author = &authors[nAuthors]
		authmap[a.name] = author
	}
	return author, nil
}

func processCommit(c *commit, a *author) {
	if strings.HasSuffix(c.message, ".") {
		c.hasEndingPeriod = true
		c.message = c.message[:len(c.message)-1]
	}
	// lowering caps improves verb detection
	c.message = strings.ToLower(c.message)
	c.user = a
}

// errSkipCommit tells program to ignore commit message
var errSkipCommit = errors.New("this commit will be ignored")

func scanNextCommit(rdr *bufio.Reader) (c commit, a author, err error) {
	var line string
	var commitLineScanned bool
	for {
		line, err = scanNextLine(rdr)
		switch {
		case !commitLineScanned && strings.HasPrefix(line, "commit"):
			commitLineScanned = true
		case strings.HasPrefix(line, "Author:"):
			a, err = parseAuthor(line[len("Author:"):])
		case strings.HasPrefix(line, "Date:"):
			c.date, err = time.Parse("Mon Jan 2 15:04:05 2006 -0700", strings.TrimSpace(line[len("Date:"):]))
		case strings.HasPrefix(line, "Merge:"):
			c.isMerge = true
		case strings.HasPrefix(line, "fatal:"):
			err = errors.New(line)
		default:
			c.message = appendMessage(c.message, line)
		}
		if err != nil {
			break
		}
		b, err := rdr.Peek(len("\ncommit"))
		if err != nil || string(b) == "\ncommit" {
			break
		}
	}
	return c, a, err
}

func scanNextLine(rdr *bufio.Reader) (string, error) {
	for {
		b, err := rdr.ReadBytes('\n')
		if err != nil {
			return "", err
		}
		if len(b) == 1 { // no text on line
			continue
		}
		return string(b[:len(b)-1]), nil
	}
}

func appendMessage(msg, toAppend string) string {
	if msg == "" {
		return strings.TrimSpace(toAppend)
	}

	return msg + " " + strings.TrimSpace(toAppend)
}

func parseAuthor(s string) (author, error) {
	mailstart := strings.Index(s, "<")
	mailend := strings.Index(s, ">")
	if mailstart < 1 || mailend < 3 {
		return author{}, errors.New("bad author line:" + s)
	}
	return author{name: strings.TrimSpace(s[:mailstart]), email: s[mailstart+1 : mailend]}, nil
}
