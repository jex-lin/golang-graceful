## Run

This is a simple example that show you the idea how to let your app shutdown or restart gracefully.

#### window 1

    $ go build && ./golang-graceful-example
    (pid: 13468) Started...
    (pid: 13468) Doing job...
    (pid: 13468) Running successfully

#### window 2

Restart with zero-down time

    $ kill -1 13468
    (pid: 13468) Forking...
    (pid: 13468) Kill self, terminating...
    (pid: 14107) Started...                         <= New process started.
    (pid: 14107) Doing job...
    (pid: 14107) Running successfully.
    (pid: 13468) Job done!
    (pid: 13468) Stop doing jobs.
    (pid: 13468) Terminated.                        <= Parent process killed.
    (pid: 13468) Parent process has been killed.


Shutdown

    $ kill 13468
    (pid: 13468) Terminating...
    (pid: 13468) Job done!
    (pid: 13468) Stop doing jobs.
    (pid: 13468) Terminated.


## Note

* Do not try to catch `SIGKILL` and `SIGSTOP` in code, there is reason : [os/signal: Prevent developers from catching SIGKILL and SIGSTOP](https://github.com/golang/go/issues/9463)
