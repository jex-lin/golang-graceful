package main

import (
	"context"
	"log"
	"net/http"
	"os"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var grace GraceHandler
var pid int

type GraceHandler struct {
	StopCh chan bool
	SyncWG sync.WaitGroup
	Http   *http.Server
}
type server struct{}

func main() {
	http_port := "3333"
	pid = syscall.Getpid()
	grace.StopCh = make(chan bool)

	log.Printf("(pid: %d) Started...\n", pid)
	log.Printf("(pid: %d) Port: %s\n", pid, http_port)

	go startHttpServer(http_port)
	go startWorker()
	handleSignals()
}

func handleSignals() {
	sig_chan := make(chan os.Signal, 1)
	signal.Notify(sig_chan, syscall.SIGINT, syscall.SIGTERM)
	log.Printf("(pid: %d) Running successfully.\n", pid)

	switch <-sig_chan {
	// Stop
	case syscall.SIGINT, // 2
		syscall.SIGTERM: // 15
		log.Printf("(pid: %d) Terminating...\n", pid)

		close(grace.StopCh) // stop worker

		// stop http server
		if grace.Http != nil {
			log.Println("Shutdown with timeout: 10 seconds")
			ctx, cancel := context.WithTimeout(context.Background(), 10*time.Second)
			defer cancel()
			if err := grace.Http.Shutdown(ctx); err != nil {
				log.Printf("Error: %v\n", err)
			} else {
				log.Printf("(pid: %d) Server stopped.\n", pid)
			}
		} else {
			log.Printf("(pid: %d) Http server didn't work.\n", pid)
		}

		grace.SyncWG.Wait()
		log.Printf("(pid: %d) Terminated.\n", pid)
		os.Exit(0)
	}
}

func startWorker() {
	// Do job
	for {
		select {
		case <-grace.StopCh:
			log.Printf("(pid: %d) Stop doing job and wait for processing job to be done.\n", pid)
			return
		default:
			go doJob()
		}
		time.Sleep(3 * time.Second)
	}
}

func doJob() {
	grace.SyncWG.Add(1)
	defer grace.SyncWG.Done()
	log.Printf("(pid: %d) Doing job ...\n", pid)
	time.Sleep(8 * time.Second)
	log.Printf("(pid: %d) Job done!\n", pid)
}

func startHttpServer(port string) {
	grace.Http = &http.Server{Addr: ":" + port}
	http.HandleFunc("/", index)
	if err := grace.Http.ListenAndServe(); err != http.ErrServerClosed {
		log.Fatal(err)
	}
}

func index(w http.ResponseWriter, r *http.Request) {
	log.Printf("(pid: %d) Doing request ...\n", pid)
	time.Sleep(8 * time.Second)
	log.Printf("(pid: %d) Request done!\n", pid)
	w.Write([]byte("Request done!\n"))
}
