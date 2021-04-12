[![Go Report Card](https://goreportcard.com/badge/github.com/soypat/gitaligned)](https://goreportcard.com/report/github.com/soypat/gitaligned)
[![go.dev reference](https://pkg.go.dev/badge/github.com/soypat/gitaligned)](https://pkg.go.dev/github.com/soypat/gitaligned)
[![codecov](https://codecov.io/gh/soypat/gitaligned/branch/main/graph/badge.svg)](https://codecov.io/gh/soypat/gitaligned/branch/main)


# gitaligned
Find out where you fall on the Open-Source Character Alignment Chart
---

Binaries available in [releases](https://github.com/soypat/gitaligned/releases).

If you prefer to install from source, run the following in your command line (requires Go)
```
go get -u github.com/soypat/gitaligned
```

## How to use (requires git)
Run `gitaligned -h` for help.

Running gitaligned in this repo:
```
gitaligned -u soypat
```

Output:
```
Author soypat is Neutral Good
Commits: 6
Accumulated:{-0.2 2}
```



## Planned

Output:
```
Steve -- Chaotic Neutral (89.9% confidence)
  82 commits
  99 % Lean towards Chaotic
  10 % Lean towards Good
```

### How it works (sort of)

For now gitaligned does some basic natural language processing using [`prose`](https://github.com/jdkato/prose) and has some ad-hoc rules based on typical git commit message mannerisms.

***Opinions*** of commit messages and their alignment

| Commit Message | Alignment |
|---|---|
| `Fixed bug`  |   Neutral Evil  |
|`Correct edge case in http response where long frames would overflow` | Lawful Good |
| `Steve's parser was really bad. Optimize and now works with extended unicode` | Chaotic Good |
| `Jacobian not singular` | True Neutral |
| `f*cking BNF` | Chaotic Evil |

"Good" commit messages seek to inform the reader of what changed and why.

"Evil" commit messages hide away what changed.


### Author alignments on this repo

| Author | Alignment |
|--------|-----------|
| soypat | Chaotic Good |
| frenata | Lawful Neutral | 