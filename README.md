# Graceful shutdown and restart

This is a simple example that show you how to let your app shutdown or restart gracefully and achieve zero-down time.




# Run example

### Window 1) Start

    $ go build && ./golang-graceful-example
    (pid: 94178) Started...
    (pid: 94178) Running successfully.

### Window 2) Request and you will get blocked for 10 seconds.

    $ curl localhost:3333

### Window 3) Restart or shutdown to see what happend.

Restart with zero-down time

    $ kill -1 94178

Shutdown

    $ kill 94178





# Issues

* Graceful shutdown http server. Closing listener will interrupt processing request.




# Note

* It doesn't work to catch signal of `SIGKILL` and `SIGSTOP`, there is the reason : [os/signal: Prevent developers from catching SIGKILL and SIGSTOP](https://github.com/golang/go/issues/9463)
