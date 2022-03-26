package utils

import (
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

func DoWithPProf(outfile string, job func()) {
	file, _ := os.OpenFile(outfile, os.O_CREATE|os.O_WRONLY, 0644)
	defer func() {
		file.Close()
		log.Printf("self pprof is sampled to %s", outfile)
	}()
	pprof.StartCPUProfile(file)
	defer pprof.StopCPUProfile()

	done := make(chan struct{})
	go func() {
		defer close(done)
		job()
	}()

	signals := make(chan os.Signal, 1)
	signal.Notify(signals, syscall.SIGINT, syscall.SIGHUP, syscall.SIGTERM)
	select {
	case sig := <-signals:
		log.Printf("interrupted by %s", sig)
	case <-done:
	}

	return
}
