package khatru

import (
	"context"
	"log"
	"net/http"
	"os"
	"sync"
	"time"

	"github.com/fasthttp/websocket"
	"github.com/nbd-wtf/go-nostr"
	"github.com/nbd-wtf/go-nostr/nip11"
)

func NewRelay() *Relay {
	rl := &Relay{
		Log: log.New(os.Stderr, "[khatru-relay] ", log.LstdFlags),

		Info: &nip11.RelayInformationDocument{
			Software:      "https://github.com/fiatjaf/khatru",
			Version:       "n/a",
			SupportedNIPs: []int{1, 11, 42, 70, 86},
		},

		upgrader: websocket.Upgrader{
			ReadBufferSize:  1024,
			WriteBufferSize: 1024,
			CheckOrigin:     func(r *http.Request) bool { return true },
		},

		clients:   make(map[*WebSocket][]listenerSpec, 100),
		listeners: make([]listener, 0, 100),

		serveMux: &http.ServeMux{},

		WriteWait:      10 * time.Second,
		PongWait:       60 * time.Second,
		PingPeriod:     30 * time.Second,
		MaxMessageSize: 512000,
	}

	return rl
}

type Relay struct {
	ServiceURL string

	// hooks that will be called at various times
	RejectEvent               []func(ctx context.Context, event *nostr.Event) (reject bool, msg string)
	OverwriteDeletionOutcome  []func(ctx context.Context, target *nostr.Event, deletion *nostr.Event) (acceptDeletion bool, msg string)
	StoreEvent                []func(ctx context.Context, event *nostr.Event) error
	DeleteEvent               []func(ctx context.Context, event *nostr.Event) error
	OnEventSaved              []func(ctx context.Context, event *nostr.Event)
	OnEphemeralEvent          []func(ctx context.Context, event *nostr.Event)
	RejectFilter              []func(ctx context.Context, filter nostr.Filter) (reject bool, msg string)
	RejectCountFilter         []func(ctx context.Context, filter nostr.Filter) (reject bool, msg string)
	OverwriteFilter           []func(ctx context.Context, filter *nostr.Filter)
	OverwriteCountFilter      []func(ctx context.Context, filter *nostr.Filter)
	QueryEvents               []func(ctx context.Context, filter nostr.Filter) (chan *nostr.Event, error)
	CountEvents               []func(ctx context.Context, filter nostr.Filter) (int64, error)
	RejectConnection          []func(r *http.Request) bool
	OnConnect                 []func(ctx context.Context)
	OnDisconnect              []func(ctx context.Context)
	OverwriteRelayInformation []func(ctx context.Context, r *http.Request, info nip11.RelayInformationDocument) nip11.RelayInformationDocument
	OverwriteResponseEvent    []func(ctx context.Context, event *nostr.Event)
	PreventBroadcast          []func(ws *WebSocket, event *nostr.Event) bool

	// these are used when this relays acts as a router
	routes                []Route
	getSubRelayFromEvent  func(*nostr.Event) *Relay // used for handling EVENTs
	getSubRelayFromFilter func(nostr.Filter) *Relay // used for handling REQs

	// setting up handlers here will enable these methods
	ManagementAPI RelayManagementAPI

	// editing info will affect the NIP-11 responses
	Info *nip11.RelayInformationDocument

	// Default logger, as set by NewServer, is a stdlib logger prefixed with "[khatru-relay] ",
	// outputting to stderr.
	Log *log.Logger

	// for establishing websockets
	upgrader websocket.Upgrader

	// keep a connection reference to all connected clients for Server.Shutdown
	// also used for keeping track of who is listening to what
	clients      map[*WebSocket][]listenerSpec
	listeners    []listener
	clientsMutex sync.Mutex

	// set this to true to support negentropy
	Negentropy bool

	// in case you call Server.Start
	Addr       string
	serveMux   *http.ServeMux
	httpServer *http.Server

	// websocket options
	WriteWait      time.Duration // Time allowed to write a message to the peer.
	PongWait       time.Duration // Time allowed to read the next pong message from the peer.
	PingPeriod     time.Duration // Send pings to peer with this period. Must be less than pongWait.
	MaxMessageSize int64         // Maximum message size allowed from peer.
}
