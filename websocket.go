package main

import (
	"context"
	"net/http"
	"sync"
	"time"

	"nhooyr.io/websocket"
)

type WebSocket struct {
	Subscribers map[*websocket.Conn]struct{}
	Mu          sync.Mutex
}

// TODO - remove this
var acceptOptions = &websocket.AcceptOptions{
	InsecureSkipVerify: true,
}

func (ws *WebSocket) AddSubscriber(w http.ResponseWriter, r *http.Request) error {
	subscriber, err := websocket.Accept(w, r, acceptOptions)
	if err != nil {
		return err
	}
	defer ws.RemoveSubscriber(subscriber)
	ctx := subscriber.CloseRead(r.Context())

	ws.Mu.Lock()
	ws.Subscribers[subscriber] = struct{}{}
	ws.Mu.Unlock()

	ws.WriteSubscribers()

	<-ctx.Done()
	return nil
}

func (ws *WebSocket) RemoveSubscriber(subscriber *websocket.Conn) {
	if subscriber == nil {
		return
	}
	subscriber.Close(websocket.StatusNormalClosure, "closed")

	ws.Mu.Lock()
	delete(ws.Subscribers, subscriber)
	ws.Mu.Unlock()
}

func (ws *WebSocket) WriteSubscribers() {
	JSON, err := tree.ToJSON()
	if err != nil {
		return
	}

	for subscriber := range ws.Subscribers {
		ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
		defer cancel()
		subscriber.Write(ctx, websocket.MessageText, JSON)
	}
}
