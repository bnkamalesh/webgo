package sse

import (
	"bytes"
	"net/http"
	"strconv"
	"time"
)

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
