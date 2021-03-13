# Gogot
### Git, but in Go

Following [this article by ThoughtBot](https://thoughtbot.com/blog/rebuilding-git-in-ruby) titled **"Rebuilding Git in Ruby"**, I made a similar program but in Go. Not really a full copy of Git, rather more like an attempt to understand how Git stores and manages file contents. The commands supported (so far) are these:

- [x] `gogot init [PATH]` - Similar to `git init [PATH]`
- [x] `gogot add [FILE1] [FILE2] [PATH]` - Similar to `git add ...`. Right now there's no support for `gogot add -A|--all`, but `gogot add .` does the trick
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
$ mv ../gogot . # Or add it to path, that's a better idea

$ touch hello.txt
$ ./gogot add .
$ ./gogot commit Initial commit

$ echo "Howdy y'all" >> hello.txt
$ ./gogot add .
$ ./gogot commit Greetings in hello.txt

$ echo "How's everyone doin'?" >> hello.txt
$ ./gogot add .
$ ./gogot commit Greetings extended...

$ ./gogot log
Logs on branch main
jqF2823Ila5NwJHzU-4K40LpaLM=    (author aldo)    Initial commit
7g7D7vGQdZ-xJ2_7Aw0MVw2eDn0=    (author aldo)    Greetings in hello.txt
DbUW6I6h3xfLlvDbqweqlDYcRuM=    (author aldo)    Greetings extended...

$ ./gogot time-machine 7g7D7vGQdZ-xJ2_7Aw0MVw2eDn0= hello.txt
Howdy y'all

```
