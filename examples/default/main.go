package main

import (
	"fmt"
	"net/http"
	"os"

	"github.com/fiatjaf/eventstore/sqlite3"
	"github.com/fiatjaf/khatru"
	"github.com/fiatjaf/khatru/policies"
)

func main() {
    relay := khatru.NewRelay()

	// NIP-11 info
	relay.Info.Name = "relay.nosto.re"
  relay.Info.PubKey = "npub1qe3e5wrvnsgpggtkytxteaqfprz0rgxr8c3l34kk3a9t7e2l3acslezefe"
  relay.Info.Contact = "info@sebastix.nl"
  relay.Info.Description = "Nostr relay optimized for NIP-B7 Blossom and NIP-5A Static Websites"
  relay.Info.Version = "0.1"

	db := sqlite3.SQLite3Backend{DatabaseURL: "./data/khatru-sqlite"}
	os.MkdirAll("./data", 0755)
	if err := db.Init(); err != nil {
		panic(err)
	}

	relay.StoreEvent = append(relay.StoreEvent, db.SaveEvent)
	relay.QueryEvents = append(relay.QueryEvents, db.QueryEvents)
	relay.CountEvents = append(relay.CountEvents, db.CountEvents)
	relay.DeleteEvent = append(relay.DeleteEvent, db.DeleteEvent)

  allowedEventKinds := []uint16{24242,10063,34128,15128,35128}
  relay.RejectEvent = append(relay.RejectEvent, policies.RestrictToSpecifiedKinds(false, allowedEventKinds...))

  // Output when there is HTTP request
  mux := relay.Router()
  // set up other http handlers
  mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
      w.Header().Set("content-type", "text/html")
      fmt.Fprintf(w, `<html><head></head><body>`)
      fmt.Fprintf(w, `<div>`)
      fmt.Fprintf(w, `<br /><br />`)
      fmt.Fprintf(w, `This relay only accepts events with kinds:`)
      fmt.Fprintf(w, `<ul>`)
      fmt.Fprintf(w, `<li><code>24242</code> (Authorization event)</li>`)
      fmt.Fprintf(w, `<li><code>10063</code> (User Blossom servers list event)</li>`)
      fmt.Fprintf(w, `<li><code>34128</code> (nsite v1: static file event)</li>`)
      fmt.Fprintf(w, `<li><code>15128</code> (nsite v2: root site manifest event)</li>`)
      fmt.Fprintf(w, `<li><code>35128</code> (nsite v2: named site manifestl event)</li>`)
      fmt.Fprintf(w, `</ul>`)
      fmt.Fprintf(w, `<br /><br />`)
      fmt.Fprintf(w, `<a href="https://github.com/Sebastix/khatru/tree/relay.nosto.re" target="https://github.com/Sebastix/khatru/tree/relay.nosto.re">https://github.com/Sebastix/khatru/tree/relay.nosto.re</a>`)
      fmt.Fprintf(w, `</div>`)
      fmt.Fprintf(w, `</body></html>`)
  })

	fmt.Println("running on :3338")
	http.ListenAndServe(":3338", relay)
}
