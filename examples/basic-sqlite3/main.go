package main

import (
	"fmt"
	"net/http"

	"github.com/fiatjaf/eventstore/sqlite3"
	"github.com/fiatjaf/khatru"
  "github.com/fiatjaf/khatru/policies"
)

func main() {
	relay := khatru.NewRelay()

	// NIP-11 info
    relay.Info.Name = "relay.kubo.watch"
    relay.Info.PubKey = "npub1kdstrkmhv0yx8pdqcf9ed8l26752gqprx68twg7qp5nsd7qtegnsr3nsze"
    relay.Info.Contact = "info@sebastix.nl"
    relay.Info.Description = ""
    relay.Info.Version = "0.0.1"

	db := sqlite3.SQLite3Backend{DatabaseURL: "/var/www/kubo/relay.kubo.watch/data/khatru-sqlite"}
	if err := db.Init(); err != nil {
		panic(err)
	}

	relay.StoreEvent = append(relay.StoreEvent, db.SaveEvent)
	relay.QueryEvents = append(relay.QueryEvents, db.QueryEvents)
	relay.CountEvents = append(relay.CountEvents, db.CountEvents)
	relay.DeleteEvent = append(relay.DeleteEvent, db.DeleteEvent)
	relay.ReplaceEvent = append(relay.ReplaceEvent, db.ReplaceEvent)

  allowedEventKinds := []uint16{0,3,5,1984,1985,1111,21,22,34235,34236,10040,30382,30383,30384,30385}
	relay.RejectEvent = append(relay.RejectEvent, policies.RestrictToSpecifiedKinds(true, allowedEventKinds[0]))

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
        fmt.Fprintf(w, `Connect your Nostr client to <code>wss://relay.kubo.watch</code>`)
        fmt.Fprintf(w, `<br /><br />`)
        fmt.Fprintf(w, `This relay only accepts events from authors: @todo`)
        fmt.Fprintf(w, `<br /><br />`)
        fmt.Fprintf(w, `This relay only accepts events with kind:`)
        fmt.Fprintf(w, `<br />`)
        fmt.Fprintf(w, `<code>0, 3, 5</code>`)
        fmt.Fprintf(w, `<br />`)
        fmt.Fprintf(w, `<code>1984, 1985</code>`)
        fmt.Fprintf(w, `<br />`)
        fmt.Fprintf(w, `<code>1111</code> (comment <a href="https://nips.nostr.com/22">NIP-22</a>)`)
        fmt.Fprintf(w, `<br />`)
        fmt.Fprintf(w, `<code>10040, 30382, 30383, 30384, 30385</code> (comment <a href="https://nips.nostr.com/85">NIP-85</a>)`)
        fmt.Fprintf(w, `<br />`)
        fmt.Fprintf(w, `<code>21</code>, <code>22</code>, <code>34235</code>, <code>34236</code> (comment <a href="https://nips.nostr.com/71">NIP-71</a>)`)
        fmt.Fprintf(w, `<br /><br />`)
        fmt.Fprintf(w, `<a href="https://github.com/nostrver-se/khatru/tree/relay.kubo.watch" target="_blank">https://github.com/nostrver-se/khatru/tree/relay.kubo.watch</a>`)
        fmt.Fprintf(w, `</div>`)
        fmt.Fprintf(w, `</body></html>`)
    })

	fmt.Println("running on :1212")
	http.ListenAndServe(":1212", relay)
}
