package main

import (
	"fmt"
	"log"
	"os"
	"os/exec"
	"os/signal"
	"sync"
	"syscall"
	"time"
)

var sigHandler SigHandler

type SigHandler struct {
	StopCh chan bool
	SyncWG sync.WaitGroup
}

func main() {
	pid := syscall.Getpid()
	fmt.Printf("(pid: %d) Started...\n", pid)

	sigHandler.StopCh = make(chan bool)

	go startWorking(pid)

	handleSignals(pid)

}

func handleSignals(pid int) {
	var sig os.Signal
	sig_chan := make(chan os.Signal, 1)
	signal.Notify(sig_chan, syscall.SIGHUP, syscall.SIGINT, syscall.SIGTERM)
	fmt.Printf("(pid: %d) Running successfully.\n", pid)

	sig = <-sig_chan
	switch sig {
	// Restart
	case syscall.SIGHUP: // 1
		fmt.Printf("(pid: %d) Forking...\n", pid)
		fork()
		fmt.Printf("(pid: %d) Kill self, terminating...\n", pid)
		close(sigHandler.StopCh) // Notify worker to stop doing new job.
		sigHandler.SyncWG.Wait() // Wait for jobs to be done.
		fmt.Printf("(pid: %d) Terminated.\n", pid)
		syscall.Kill(syscall.Getpid(), syscall.SIGTERM)
		fmt.Printf("(pid: %d) Parent process has been killed.\n", pid)

	// Stop
	case syscall.SIGINT, // 2
		syscall.SIGTERM: // 15
		fmt.Printf("(pid: %d) Terminating...\n", pid)
		close(sigHandler.StopCh)
		sigHandler.SyncWG.Wait()
		fmt.Printf("(pid: %d) Terminated.\n", pid)
		os.Exit(0)
	}
}

func startWorking(pid int) {
	for {
		select {
		case <-sigHandler.StopCh:
			fmt.Printf("(pid: %d) Stop receiving jobs.\n", pid)
			return
		default:
			doJob(pid)
		}

	}
}

func doJob(pid int) {
	sigHandler.SyncWG.Add(1)
	defer sigHandler.SyncWG.Done()

	fmt.Printf("(pid: %d) Doing job...\n", pid)
	time.Sleep(3 * time.Second)
	fmt.Printf("(pid: %d) Job done!\n", pid)
}

func fork() {
	path := "./golang-graceful-example"
	cmd := exec.Command(path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	err := cmd.Start()
	if err != nil {
		log.Fatalf("Failed to fork process, error: %v\n", err)
	}
}
