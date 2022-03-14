// Package sse implements Server-Sent Events(SSE)
// This extension is compliant with any net/http implementation, and is not limited to WebGo.
package sse

import (
	"bytes"
	"net/http"
	"strconv"
	"sync"
	"sync/atomic"
	"time"
)

type SSE struct {
	Clients      sync.Map
	clientsCount atomic.Value
	// ClientIDHeader is the HTTP request header in which the client ID is set. Default is `sse-clientid`
	ClientIDHeader string
	// UnsupportedMessage is used to send the error response to client if the
	// server doesn't support SSE
	UnsupportedMessage func(http.ResponseWriter, *http.Request) error

	// OnCreateClient is a hook, for when a client is added to the active clients
	OnCreateClient func(clientID string, count int)

	// OnRemoveClient is a hook, for when a client is removed from the active clients
	OnRemoveClient func(clientID string, count int)
}

// Handler returns an error rather than being directly used as an http.HandlerFunc,
// to let the user handle error. e.g. if the error has to be logged
func (sse *SSE) Handler(w http.ResponseWriter, r *http.Request) error {
	flusher, hasFlusher := w.(http.Flusher)
	if !hasFlusher {
		return sse.UnsupportedMessage(w, r)
	}

	header := w.Header()
	header.Set("Content-Type", "text/event-stream")
	header.Set("Connection", "keep-alive")
	header.Set("X-Accel-Buffering", "no")
	w.WriteHeader(http.StatusOK)

	clientID := r.Header.Get(sse.ClientIDHeader)
	msg, _ := sse.ClientMessageChan(clientID)
	defer sse.RemoveClientMessageChan(clientID)

	ctx := r.Context()
	for {
		select {

		case payload, ok := <-msg:
			if !ok {
				return nil
			}

			_, err := w.Write(payload.Bytes())
			if err != nil {
				close(msg)
				return err
			}

		case <-ctx.Done():
			{
				return ctx.Err()
			}
		}

		flusher.Flush()
	}
}

// HandlerFunc is a convenience function which can be directly used with net/http implementations.
// Important: You cannot handle any error returned by the Handler
func (sse *SSE) HandlerFunc(w http.ResponseWriter, r *http.Request) {
	_ = sse.Handler(w, r)
}

// ClientMessageChan returns a message channel to stream data to a client
// The boolean value is `true` if the client didn't exist before
func (sse *SSE) ClientMessageChan(clientID string) (chan *Message, bool) {
	msg, ok := sse.Clients.Load(clientID)
	if !ok {
		msg = make(chan *Message)
		sse.Clients.Store(clientID, msg)
		count := sse.clientsCount.Load().(int)
		count++
		sse.clientsCount.Store(count)
		sse.OnCreateClient(clientID, count)
	}

	return msg.(chan *Message), !ok
}

// RemoveClientMessageChan removes the channel from clients map given a clientID
func (sse *SSE) RemoveClientMessageChan(clientID string) {
	sse.Clients.Delete(clientID)
	count := sse.clientsCount.Load().(int)
	count--
	sse.clientsCount.Store(count)

	sse.OnRemoveClient(clientID, count)
}

// Broadcast sends the message to all active clients
func (sse *SSE) Broadcast(msg Message) {
	sse.Clients.Range(func(key, value interface{}) bool {
		mchan, _ := value.(chan *Message)
		mchan <- &msg
		return true
	})
}

func (sse *SSE) ActiveClients() int {
	return sse.clientsCount.Load().(int)
}

// Message represents a valid SSE message
// ref: https://developer.mozilla.org/en-US/docs/Web/API/Server-sent_events/Using_server-sent_events
type Message struct {
	// Event is a string identifying the type of event described. If this is specified, an event will be dispatched on the browser to the listener for the specified event name; the website source code should use addEventListener() to listen for named events. The onmessage handler is called if no event name is specified for a message.
	Event string

	// Data field for the message. When the EventSource receives multiple consecutive lines that begin with data:, it concatenates them, inserting a newline character between each one. Trailing newlines are removed.
	Data string

	// ID to set the EventSource object's last event ID value.
	ID string

	// Retry is the reconnection time. If the connection to the server is lost, the browser will wait for the specified time before attempting to reconnect. This must be an integer, specifying the reconnection time in milliseconds. If a non-integer value is specified, the field is ignored.
	Retry time.Duration
}

func (m *Message) Bytes() []byte {
	// The event stream is a simple stream of text data which must be encoded using UTF-8.
	// Messages in the event stream are separated by a pair of newline characters.
	// A colon as the first character of a line is in essence a comment, and is ignored.

	buff := bytes.NewBufferString("")
	if m.Event != "" {
		buff.WriteString("event:" + m.Event + "\n")
	}
	if m.ID != "" {
		buff.WriteString("id:" + m.ID + "\n")
	}
	if m.Data != "" {
		buff.WriteString("data:" + m.Data + "\n")
	}
	if m.Retry != 0 {
		buff.WriteString("retry:" + strconv.Itoa(int(m.Retry.Milliseconds())) + "\n")
	}
	buff.WriteString("\n")
	return buff.Bytes()
}

func DefaultUnsupportedMessageHandler(w http.ResponseWriter, r *http.Request) error {
	w.WriteHeader(http.StatusNotImplemented)
	_, err := w.Write([]byte("Streaming not supported"))
	return err
}

func DefaultHook(clientID string, count int) {}

func New() *SSE {
	clientsCount := atomic.Value{}
	clientsCount.Store(int(0))

	s := &SSE{
		Clients:            sync.Map{},
		clientsCount:       clientsCount,
		ClientIDHeader:     "sse-clientid",
		UnsupportedMessage: DefaultUnsupportedMessageHandler,
		OnRemoveClient:     DefaultHook,
		OnCreateClient:     DefaultHook,
	}

	return s
}
