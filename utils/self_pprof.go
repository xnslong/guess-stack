package utils

import (
	"io"
	"log"
	"os"
	"os/signal"
	"runtime/pprof"
	"syscall"
)

func DoWithPProf(outfile string, job func()) {
	_ = WriteToFile(outfile, func(writer io.Writer) error {
		_ = pprof.StartCPUProfile(writer)
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

		log.Printf("salf-pprof is written to: %s", outfile)
		return nil
	})
}
