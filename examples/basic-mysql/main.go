package main

import (
	"fmt"
	"net/http"

	"github.com/fiatjaf/eventstore/mysql"
	"github.com/fiatjaf/khatru"
)

func main() {
	relay := khatru.NewRelay()

	// NIP-11 info
    relay.Info.Name = "khatru.nostrver.se"
    relay.Info.PubKey = "npub1qe3e5wrvnsgpggtkytxteaqfprz0rgxr8c3l34kk3a9t7e2l3acslezefe"
    relay.Info.Contact = "info@sebastix.nl"
    relay.Info.Description = "Custom relay build with Khatru (mysql)"
    relay.Info.Version = "0.0.1"

    db := mysql.MySQLBackend{DatabaseURL: '...'}
    // @TODO find out how to connect to a mysql database which is running on the server
    if err := db.Init(); err != nil {
        panic(err)
    }

    fmt.Println("running on :3334")
    http.ListenAndServe(":3334", relay)

}