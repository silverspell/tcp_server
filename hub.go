package tcp_server

import (
	"fmt"
	"os"
	"strconv"
	"time"
)

type Hub struct {
	Clients map[string]*Client
	cleaner time.Timer
}

func NewHub() *Hub {
	h := &Hub{
		Clients: map[string]*Client{},
		cleaner: *time.NewTimer(1 * time.Minute),
	}
	go h.clearConnections()
	return h
}

func (hub *Hub) clearConnections() {
	<-hub.cleaner.C
	fmt.Println("Starting cleaner")
	var maxTimeDiffBetweenMessages int64

	envValue, ok := os.LookupEnv("MAX_TIME_DIFF_BETWEEN_MESSAGES")
	if ok {
		maxTimeDiffBetweenMessages, _ = strconv.ParseInt(envValue, 10, 64)
	}
	now := time.Now().Unix()

	for key, client := range hub.Clients {
		if now-client.lastSeen >= maxTimeDiffBetweenMessages {
			client.conn.Close()
			delete(hub.Clients, key)
			fmt.Printf("%s disconnected\n", key)
		}
	}
	hub.cleaner.Reset(1 * time.Minute)
	fmt.Println("Cleaner reset")
	go hub.clearConnections()
}
