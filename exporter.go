package main

import (
	"context"
	"flag"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/jonboulle/clockwork"
	"github.com/prometheus/client_golang/prometheus/promhttp"
)

var (
	listen = flag.String("listen", ":9099", "Exporter listen address")
	target = flag.String("target", "", "IP address of your Nokia Gateway device")

	scrapeFrequency = flag.Duration("scrape_frequency", 10*time.Second, "How frequently to get the status from the Gateway")
)

func main() {
	log.SetOutput(os.Stdout)
	log.SetFlags(log.Lshortfile)

	flag.Parse()
	if *target == "" {
		log.Fatalf("the required -target flag is currently unset, exiting.")
	}

	c := Collector{
		*target,
	}

	ctx, cancel := context.WithCancel(context.Background())
	defer cancel()

	go c.ScrapeEvery(ctx, clockwork.NewRealClock(), *scrapeFrequency)

	http.Handle("/metrics", promhttp.Handler())

	log.Printf("listening for connections at %v", *listen)
	if err := http.ListenAndServe(*listen, nil); err != nil {
		log.Fatalf("ListenAndServe at %v failed: %v", *listen, err)
	}
}
