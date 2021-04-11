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
	user            *author
	hasEndingPeriod bool
	date            time.Time
	message         string
}
type author struct {
	commits int
	alignment
	accumulator alignment
	name, email string
}

func (a author) Stats() string {
	return fmt.Sprintf("Author %v is %v\nCommits: %d\nAccumulated:%0.1g\n",
		a.name, a.alignment.Format(), a.commits, a.accumulator)
}

// ScanCWD Scans .git in current working directory using git
// command. Scans up to maxCommit messages.
func ScanCWD(branch string) ([]commit, []author, error) {
	if branch == "" {
		branch = "--all"
	}
	cmd := exec.Command("git", "log", branch, "-n", strconv.Itoa(maxCommits))
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
	authors = make([]author, 0, maxAuthors)
	authmap := make(map[string]*author)
	var b []byte
	c := commit{}
	counter := 0
	var skipFlag bool
	for {
		if counter >= maxCommits { // || (rdr.Buffered() == 0 && counter != 0) {
			break
		}
		b, err = rdr.ReadBytes('\n')
		if err != nil {
			break
		}
		if len(b) == 1 { // no text on line
			continue
		}
		line := string(b[:len(b)-1])
		switch {
		case strings.HasPrefix(line, "commit"):
			skipFlag = false
			if c.user != nil {
				if strings.HasSuffix(c.message, ".") {
					c.hasEndingPeriod = true
					c.message = c.message[:len(c.message)-1]
				}
				// lowering caps improves verb detection
				c.message = strings.ToLower(c.message)
				commits = append(commits, c)
				counter++
			}
			c = commit{}
		case strings.HasPrefix(line, "Author:"):
			a, err := parseAuthor(line[len("Author:"):])
			if err != nil {
				break
			}
			if username != "" && username != a.name {
				skipFlag = true
				c = commit{}
				continue
			}
			author, ok := authmap[a.name]
			if !ok {
				if len(authors) == maxAuthors {
					skipFlag = true
					continue
				}
				authors = append(authors, a)
				author = &authors[len(authors)-1]
				authmap[a.name] = author
			}
			c.user = author
		case strings.HasPrefix(line, "Date:"):
			if skipFlag {
				continue
			}
			c.date, err = time.Parse("Mon Jan 2 15:04:05 2006 -0700", strings.TrimSpace(line[len("Date:"):]))
			if err != nil {
				return commits, authors, err
			}
		case strings.HasPrefix(line, "Merge:"):
			continue
		case strings.HasPrefix(line, "fatal:"):
			err = errors.New(line)
		default:
			if skipFlag {
				continue
			}
			if c.message != "" {
				c.message += " "
			}
			c.message += strings.TrimSpace(line)
		}
		if err != nil {
			break
		}
	}

	if (err == nil || err == io.EOF) && c.user != nil {
		commits = append(commits, c)
	}

	return commits, authors, err
}

func parseAuthor(s string) (author, error) {
	mailstart := strings.Index(s, "<")
	mailend := strings.Index(s, ">")
	if mailstart < 1 || mailend < 3 {
		return author{}, errors.New("bad author line:" + s)
	}
	return author{name: strings.TrimSpace(s[:mailstart]), email: s[mailstart+1 : mailend]}, nil
}

// repReader spams (writes repeatedly) a string to a reader
type repReader string

func (r repReader) Read(b []byte) (int, error) {
	if len(r) < 1 {
		return 0, errors.New("bad repReader")
	}
	for i := range b {
		b[i] = byte((r)[i%len(r)])
	}
	return len(b), nil
}
