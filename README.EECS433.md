Readme for EECS 433 Project Submission
======================================

This is the submission for Stephen Brennan's project, An XPath Query Evaluator
for Filesystems. The report PDF is located in the file:

    Brennan.Report.XPathFilesystem.pdf

You may build and run DPath on Mac or Linux (assuming you have Go installed)
with the following sequence of commands:

0. `export GPATH=$HOME/go` if you do not have a Go environment set up
1. `go get` will install dependencies to your Go environment
2. `go build` will compile DPath. Executable will be simply named `dpath`
3. `go test` will run unit tests, if you care to run them

A pre-compiled version for 64 bit Linux is included in this submission as
`dpath.linux`.

Files
-----

The file `SYNTAX.md` describes all major features of DPath and gives basic usage
instructions. The file `GUIDE.md` contains a small guide to the code. Both of
these resources have some overlap with the information given in the paper.

The source code of DPath is provided in the following files:

    dpath.nex
    dpath.y
    tree.go
    item.go
    sequence.go
    axis.go
    lib.go
    util.go
    error.go
    main.go

Additionally, the following files contain tests for language constructs:

    eval_test.go
    lexer_test.go
    lib_test.go
    parser_test.go
    testutil.go

The following files are Go code which is *generated*:
- `dpath.nn.go` is a lexer generated from `dpath.nex`

    You are welcome to regenerate this file with the commands:

        go get github.com/blynn/nex
        $GOPATH/bin/nex dpath.nex


- `y.go` is a parser generated from `dpath.y`

    You are welcome to regenerate this file with the command:

        go tool yacc dpath.y
