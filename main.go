package main

import (
	"fmt"
	"log"
	"net"
	"net/http"
	"os"
	"os/exec"
	"os/signal"
	"path/filepath"
	"sync"
	"syscall"
	"time"
)

var sigHandler SigHandler
var netListener net.Listener
var pid int

type SigHandler struct {
	StopCh chan bool
	SyncWG sync.WaitGroup
}

func main() {
	pid = syscall.Getpid()
	fmt.Printf("(pid: %d) Started...\n", pid)
	sigHandler.StopCh = make(chan bool)

	go startWorking()
	handleSignals()
}

func handleSignals() {
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

func startWorking() {
	// API
	go startHttpServer()

	// Do job
	for {
		select {
		case <-sigHandler.StopCh:
			fmt.Printf("(pid: %d) Stop doing job and wait for processing job to be done.\n", pid)
			return
		default:
			sigHandler.SyncWG.Add(1)
			go doJob()
		}
		time.Sleep(3 * time.Second)
	}
}

func doJob() {
	defer sigHandler.SyncWG.Done()
	time.Sleep(5 * time.Second)
	fmt.Printf("(pid: %d) Job done!\n", pid)
}

func startHttpServer() {
	netListener = getListener()
	http.HandleFunc("/", index)
	log.Fatal(http.Serve(netListener, nil))
}

func index(w http.ResponseWriter, r *http.Request) {
	sigHandler.SyncWG.Add(1)
	defer sigHandler.SyncWG.Done()
	fmt.Printf("(pid: %d) Request received.\n", pid)
	time.Sleep(10 * time.Second)
	w.Write([]byte("Time: " + time.Now().Format(time.RFC1123)))
}

func getListener() (l net.Listener) {
	var err error
	l, err = net.Listen("tcp", ":3333")
	if err != nil {
		// Child process will come in here because the port is already in use. So get socket copy from parent process to reuse it.
		f := os.NewFile(3, "") // 0, 1, 2 is preserved for standard input, output and error, so started with 3.
		l, err = net.FileListener(f)
		if err != nil {
			fmt.Printf("(pid: %d) Failed to inherit socket file from parent process. err: %v\n", pid, err)
			os.Exit(1)
		}
		fmt.Printf("(pid: %d) socket file inherited from parent process.\n", pid)
	}
	return
}

func fork() {
	// Get duplicate
	tl := netListener.(*net.TCPListener)
	file, _ := tl.File()

	// Current app binary file
	bin_path, err := filepath.Abs(os.Args[0])
	if err != nil {
		fmt.Println("Failed to get absolute path of current binary.")
	}

	cmd := exec.Command(bin_path)
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.ExtraFiles = []*os.File{file}
	err = cmd.Start()
	if err != nil {
		log.Fatalf("Failed to fork process, error: %v\n", err)
	}
}
