# gitaligned
Find where you fall on the Open-Source Character Alignment Chart

*This is a work in progress*

***Opinions*** of commit messages and their alignment

| Commit Message | Alignment |
|---|---|
| `Fixed bug`  |   Neutral Evil  |
|`Correct edge case in http response where long frames would overflow` | Lawful Good |
| `Steve's parser was really bad. Optimize and now works with extended unicode` | Chaotic Good |
| `Jacobian not singular` | True Neutral |
| `f*cking BNF` | Chaotic Evil |


## How to use
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
