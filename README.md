# Gogot
### Git, but in Go

There was an attempt to write a Git clone in Go. This is the attempt:

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
Logs on branch master
jqF2823Ila5NwJHzU-4K40LpaLM=    (author aldo)    Initial commit
7g7D7vGQdZ-xJ2_7Aw0MVw2eDn0=    (author aldo)    Greetings in hello.txt
DbUW6I6h3xfLlvDbqweqlDYcRuM=    (author aldo)    Greetings extended...
$ ./gogot time-machine 7g7D7vGQdZ-xJ2_7Aw0MVw2eDn0= hello.txt
Howdy y'all

```
