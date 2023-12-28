package main

import (
	http "net/http"

	"github.com/prometheus/client_golang/prometheus/promhttp"
)

func init() {
	http.Handle("/metrics", promhttp.Handler())
}

func main() {
	err := http.ListenAndServe(":5502", nil)
	if err != nil {
		panic(err)
	}
}
