package main

import (
	"bufio"
	"errors"
	"io"
	"os/exec"
	"strings"
	"time"
)

type commit struct {
	alignment
	user            *author
	hasEndingPeriod bool
	date            time.Time

	message string
}
type author struct {
	name, email string
}

// ScanCWD Scans .git in current working directory using git
// command. obtains all commit messages.
func ScanCWD() ([]commit, []author, error) {
	cmd := exec.Command("git", "log", "--all")
	reader, writer := io.Pipe()
	cmd.Stdout = writer
	go func() {
		cmd.Run()
		writer.Close()
	}()
	commits, authors, err := GitLogScan(reader)
	if err == io.EOF {
		err = nil
	}
	return commits, authors, err
}

// GitLogScan reads git log results and generates commits
func GitLogScan(r io.Reader) (commits []commit, authors []author, err error) {
	rdr := bufio.NewReader(r)
	commits = make([]commit, 0, 100)
	authors = make([]author, 0, 20)
	authmap := make(map[string]*author)
	var b []byte
	c := commit{}
	for {
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
			if c.user != nil {
				if strings.HasSuffix(c.message, ".") {
					c.hasEndingPeriod = true
					c.message = c.message[:len(c.message)-1]
				}
				// lowering caps improves verb detection
				c.message = strings.ToLower(c.message)
				commits = append(commits, c)
			}
			c = commit{}
		case strings.HasPrefix(line, "Author:"):
			a, err := parseAuthor(line[len("Author:"):])
			if err != nil {
				break
			}
			author, ok := authmap[a.name]
			if !ok {
				authors = append(authors, a)
				author = &authors[len(authors)-1]
				authmap[a.name] = author
			}
			c.user = author
		case strings.HasPrefix(line, "Date:"):
			c.date, err = time.Parse("Mon Jan 2 15:04:05 2006 -0700", strings.TrimSpace(line[len("Date:"):]))
			if err != nil {
				return commits, authors, err
			}
		default:
			c.message += strings.TrimSpace(line)
		}
	}

	if err == nil || err == io.EOF && c.user != nil {
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
