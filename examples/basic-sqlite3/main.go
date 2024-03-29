package main

import (
	"fmt"
	//"context"
	"net/http"
	"os"
	//"slices"

	"github.com/fiatjaf/eventstore/sqlite3"
	"github.com/fiatjaf/khatru"
	"github.com/fiatjaf/khatru/policies"
	//"github.com/nbd-wtf/go-nostr"
)

func main() {
    relay := khatru.NewRelay()

	// NIP-11 info
	relay.Info.Name = "khatru.nostrver.se"
    relay.Info.PubKey = "npub1qe3e5wrvnsgpggtkytxteaqfprz0rgxr8c3l34kk3a9t7e2l3acslezefe"
    relay.Info.Contact = "info@sebastix.nl"
    relay.Info.Description = "Custom relay build with Khatru"
    relay.Info.Version = "0.0.1"

	db := sqlite3.SQLite3Backend{DatabaseURL: "./data/khatru-sqlite"}
	os.MkdirAll("./data", 0755)
	if err := db.Init(); err != nil {
		panic(err)
	}

	relay.StoreEvent = append(relay.StoreEvent, db.SaveEvent)
	relay.QueryEvents = append(relay.QueryEvents, db.QueryEvents)
	relay.CountEvents = append(relay.CountEvents, db.CountEvents)
	relay.DeleteEvent = append(relay.DeleteEvent, db.DeleteEvent)

    allowedEventKinds := []uint16{37515, 13811}
	relay.RejectEvent = append(relay.RejectEvent, policies.RestrictToSpecifiedKinds(allowedEventKinds[0], allowedEventKinds[1]))

    // Custom policy
    //relay.RejectEvent = append(relay.RejectEvent,
    //    // We only accept events with kind 37515, 13811 so we put them in an array
    //    func(ctx context.Context, event *nostr.Event) (reject bool, msg string) {
    //        fmt.Printf("%T: %d \n", event.Kind, event.Kind)
    //        slices.Sort(allowedEventKinds)
    //        n, found := slices.BinarySearch(allowedEventKinds, uint16(event.Kind))
    //        fmt.Println(n, found)
    //        if found {
    //            return false, ""
    //        }
    //        return true, "This event kind not allowed on this relay"
    //    },
    //)

    // Output when there is HTTP request
    mux := relay.Router()
    // set up other http handlers
    mux.HandleFunc("/", func(w http.ResponseWriter, r *http.Request) {
        w.Header().Set("content-type", "text/html")
        fmt.Fprintf(w, `<html><head></head><body>`)
        fmt.Fprintf(w, `<div style="text-align: center;">`)
        fmt.Fprintf(w, `Connect your Nostr client to <code>wss://khatru.nostrver.se</code>`)
        fmt.Fprintf(w, `<br /><br />`)
        fmt.Fprintf(w, `This relay only accepts events with kind <code>37515</code> (places) and <code>13811</code> (check-ins)`)
        fmt.Fprintf(w, `<br /><br />`)
        fmt.Fprintf(w, `<a href="https://github.com/Sebastix/khatru" target="https://github.com/Sebastix/khatru">https://github.com/Sebastix/khatru</a>`)
        fmt.Fprintf(w, `</div>`)
        fmt.Fprintf(w, `</body></html>`)
    })

	fmt.Println("running on :3334")
	http.ListenAndServe(":3334", relay)
}
