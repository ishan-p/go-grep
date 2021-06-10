## Grep like CLI in go
## 
##
#### Problem statement

Write a command line program that implements Unix command `grep` like functionality.

#### Features required and status

- [x] Ability to search for a string in a file

```
$ ./go-grep "search_string" filename.txt
I found the search_string in the file.
```

- [-] Ability to search for a string from standard input

```
$ ./go-grep foo
bar
barbazfoo
Foobar
food
^D
```

output -

```
barbazfoo
food
```

- [-] Ability to write output to a file instead of a standard out.

```
$ ./go-grep lorem loreipsum.txt -o out.txt
```

should create an out.txt file with the output from `go-grep`. for example,

```
$ cat out.txt
lorem ipsum
a dummy text usually contains lorem ipsum
```

- [x] Ability to search for a string recursively in any of the files in a given directory. When searching in multiple files, the output should indicate the file name and all the output from one file should be grouped together in the final output. (in other words, output from two files shouldn't be interleaved in the final output being printed)

```
$ ./go-grep "test" tests
tests/test1.txt:this is a test file
tests/test1.txt:one can test a program by running test cases
tests/inner/test2.txt:this file contains a test line
```

- [x] Package test

### Instruction to install and use

```
$ make
$ ./go-grep search_string input_file
```