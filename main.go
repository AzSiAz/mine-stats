package main

import (
	log "github.com/sirupsen/logrus"
)

func main() {
	store, err := NewStore("db.storm")
	if err != nil {
		log.
			WithError(err).
			Fatal("Error opening DB")
	}
	defer store.Close()
	log.Info("Done opening database")
}
