package main

import (
	"flag"
	"fmt"
	"log"
	"os"
	"os/signal"
	"path/filepath"
	"syscall"

	"github.com/LYH2263/go-geoalert"
	"github.com/LYH2263/go-geoalert/internal/clock"
	"github.com/LYH2263/go-geoalert/internal/httpapi"
)

func main() {
	addr := flag.String("addr", ":8115", "listen address")
	web := flag.String("web", "web", "static web directory")
	data := flag.String("data", "data", "data directory")
	flag.Parse()

	_ = os.MkdirAll(*data, 0o755)
	eng, err := geoalert.New(geoalert.Options{
		Clock:        clock.Real{},
		PersistPath:  filepath.Join(*data, "fences.json"),
		AuditPath:    filepath.Join(*data, "audit.log"),
		AlertLogPath: filepath.Join(*data, "alerts.log"),
	})
	if err != nil {
		log.Fatal(err)
	}
	defer eng.Close()
	_ = eng.LoadPersisted()

	srv := httpapi.New(eng, *web)
	go func() {
		log.Printf("geod listening on %s", *addr)
		if err := srv.ListenAndServe(*addr); err != nil {
			log.Printf("server stopped: %v", err)
		}
	}()

	ch := make(chan os.Signal, 1)
	signal.Notify(ch, os.Interrupt, syscall.SIGTERM)
	<-ch
	fmt.Println("shutting down")
}
