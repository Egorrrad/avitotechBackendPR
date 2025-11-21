package main

import (
	"log"

	"github.com/Egorrrad/avitotechBackendPR/internal/app/prservice"
)

func main() {
	conf := prservice.NewConfig()

	s := prservice.New(conf)
	if err := s.Start(); err != nil {
		log.Fatal(err)
	}
}
