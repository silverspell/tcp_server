package tcp_server

import (
	"os"
	"strconv"
	"time"
)

type Hub struct {
	Clients map[string]*Client
}

func clearConnections(hub *Hub) {
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
		}
	}
}
