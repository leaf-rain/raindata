package main

import (
	_ "embed"
	"github.com/leaf-rain/raindata/common/recuperate"
	"log"
	"net/http"
	_ "net/http/pprof"
	"os"
	"os/signal"
	"syscall"
)

//go:generate wire

func main() {
	recuperate.GoSafe(func() {
		err := http.ListenAndServe(":6060", nil)
		if err != nil {
			log.Fatal("rpprof http.ListenAndServe failed", err)
		}
	})
	adapter, err := Initialize()
	if err != nil {
		panic(err)
	}
	adapter.Run()
	listenToSystemSignals(adapter.Close)
}

func listenToSystemSignals(quitFunc func()) {
	signalChan := make(chan os.Signal, 1)
	signal.Notify(signalChan, syscall.SIGINT, syscall.SIGTERM, syscall.SIGQUIT, syscall.SIGKILL, os.Kill)
	sig := <-signalChan
	log.Printf("catch signal:%+v\n", sig)
	quitFunc()
	log.Printf("server exiting.\n")
}
