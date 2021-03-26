# Gogot
### Git, but in Go

[![Build Status](https://travis-ci.com/aziflaj/gogot.svg?branch=main)](https://travis-ci.com/aziflaj/gogot)

Inspired by [this article by ThoughtBot](https://thoughtbot.com/blog/rebuilding-git-in-ruby) titled **"Rebuilding Git in Ruby"**, I made a similar program but in Go. Not really a full copy of Git, rather more like an attempt to understand how Git stores and manages file contents. The features supported (so far) are these:

- [x] `.gogotignore` - Similar to `.gitignore`
- [x] `gogot init [PATH]` - Similar to `git init [PATH]`, it initializes a repo with a default `main` branch
- [x] `gogot add [PATH1] [PATH2] ...` - Similar to `git add ...`. Right now there's no support for `gogot add -A|--all`, but `gogot add .` does the trick
- [x] `gogot commit [MESSAGE]` - Similar to `git commit -m "[MESSAGE]"`; gogot doesn't require quotes
- [x] `gogot log` - Similar to `git log --oneline`; more condensed and quicker to write
- [x] `gogot time-machine [COMMIT-ID] [FILE-PATH]`- I don't know the Git equivalent of this, but it prints the content of a given file in the specified commit

- [ ] `gogot status` - Similar to `git status`, showing different sections for not indexed files and for the files changed after indexed
- [ ] `gogot branch [NEW-BRANCH-NAME]` - Similar to `git branch -m NEW-BRANCH-NAME`
- [ ] `gogot checkout [BRANCH]` - Similar to Git's checkout for changing branches
- [ ] `gogot rollback [COMMIT-ID] [FILE-PATH:optional]` - I don't know the Git equivalent for this, but it should revert files (or just `FILE-PATH`, if it's present) as they were in `COMMIT-ID`. 

Here's a sample run:

```bash 
$ go build
$ ./gogot init kewl-projekt # You can also do `gogot init .`
Initalizing new Gogot repo
Gogot repo initialized in kewl-projekt/.gogot

$ cd kewl-projekt
$ touch hello.txt
$ ../gogot add .
$ ../gogot commit Initial commit

$ echo "Howdy y'all" >> hello.txt
$ ../gogot add .
$ ../gogot commit Greetings in hello.txt

$ echo "How's everyone doin'?" >> hello.txt
$ ../gogot add .
$ ../gogot commit Greetings extended...

$ ../gogot log
Logs on branch main
-jr-up5Sz32qNHjlxn_qsj367vc=    (author aldo)    Initial commit
S7CM7poOCDq6wSTqNL9g0uOTXrQ=    (author aldo)    Greetings in hello.txt
jYQ5y4BvBhv_n8FDSi0o0JHcsWk=    (author aldo)    Greetings extended...

$ ../gogot status
On branch main

Files not added to index:
    (use "gogot add <path>") to include in the commit

Untracked files:
    (use "gogot add <path>") to include in the commit
nothing to commit, working tree clean


$ touch SchrodingersCat
$ echo -n -e \\x61\\x6c\\x69\\x76\\x65 > SchrodingersCat
$ echo "Me myself, I'm doin' fine!" >> hello.txt
$ touch answers
$ echo 42 >> answers
$ ../gogot status
On branch main

Files not added to index:
    (use "gogot add <path>") to include in the commit
	./hello.txt

Untracked files:
    (use "gogot add <path>") to include in the commit
	./SchrodingersCat
	./answers

$ ../gogot add ./answers
$  ../gogot status
On branch main
Files to be committed:
	./answers

Files not added to index:
    (use "gogot add <path>") to include in the commit
	./hello.txt

Untracked files:
    (use "gogot add <path>") to include in the commit
	./SchrodingersCat

$ ../gogot add .
$ ../gogot status
On branch main
Files to be committed:
	./answers
	./SchrodingersCat
	./hello.txt

$ ../gogot commit A bunch of added files
$ ../gogot status
On branch main

nothing to commit, working tree clean

$ ../gogot log
Logs on branch main
-jr-up5Sz32qNHjlxn_qsj367vc=    (author aldo)    Initial commit
S7CM7poOCDq6wSTqNL9g0uOTXrQ=    (author aldo)    Greetings in hello.txt
jYQ5y4BvBhv_n8FDSi0o0JHcsWk=    (author aldo)    Greetings extended...
MieNw-fi1H0XbdhOLBmrd7cwmP4=    (author aldo)    A bunch of added files

$ ../gogot time-machine jYQ5y4BvBhv_n8FDSi0o0JHcsWk= ./hello.txt
Howdy y'all

$ cat hello.txt
Howdy y'all
How's everyone doin'?

```
