# Introduction

This is a simple example for graceful shutdown.

Requirements: go version 1.8

# Try it!

1. window 1) `go run main.go`
1. window 2) `curl 127.0.0.1:3333`
1. window 1) `ctrl`+`c` (It will block and wait for all the jobs to be done.)
1. window 2) `curl 127.0.0.1:3333` (It won't work as we expect.)


# Note

* If it doesn't work catching signals of `SIGKILL` and `SIGSTOP`, there is the reason : [os/signal: Prevent developers from catching SIGKILL and SIGSTOP](https://github.com/golang/go/issues/9463)


