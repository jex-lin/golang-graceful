# Graceful shutdown and restart

This is a simple example that show you how to let your app shutdown or restart gracefully.

It also demonstrate how to reuse the socket from parent process,
but I don't deal with http server shutdown well, please ignore that part.



# Run example

### Window 1) Start

    $ go build && ./golang-graceful-example
    (pid: 94178) Started...
    (pid: 94178) Running successfully.

### Window 2) Restart or shutdown to see the result.

Restart with zero-down time (send SIGINT)

    $ kill -1 94178

Shutdown (send SIGTERM)

    $ kill 94178





# Issues

* Graceful shutdown http server. Closing listener will interrupt processing request.




# Note

* It doesn't work to catch signal of `SIGKILL` and `SIGSTOP`, there is the reason : [os/signal: Prevent developers from catching SIGKILL and SIGSTOP](https://github.com/golang/go/issues/9463)


