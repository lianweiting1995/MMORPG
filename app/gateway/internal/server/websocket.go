package server

import (
	"net/http"
	"sync"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
)

type Connection struct {
	ws       *websocket.Conn
	sendChan chan []byte
	playerId string
	zoneID   int
}

type Gateway struct {
	connections *sync.Map          // playerID -> *Connection
	zoneBuckets [256]*sync.Map     // 按区域分桶
	upgrader    websocket.Upgrader // 连接标记
	log         log.Logger         // 日志记录
}

func InitGateway(log log.Logger) *Gateway {
	g := &Gateway{
		upgrader: websocket.Upgrader{
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
			EnableCompression: true,
		},
		log:         log,
		connections: &sync.Map{},
	}

	return g
}

func (g *Gateway) HandleConnection(w http.ResponseWriter, r *http.Request) {
	conn, err := g.upgrader.Upgrade(w, r, nil)
	if err != nil {
		panic(err)
	}

	c := &Connection{
		ws:       conn,
		sendChan: make(chan []byte, 100),
		playerId: uuid.New().String(),
	}
	g.connections.Store(c.playerId, c)

	go g.readPump(c)
	go g.writePump(c)
}

func (g *Gateway) writePump(c *Connection) {
	ticker := time.NewTicker(15 * time.Second)
	defer ticker.Stop()

	for {
		select {
		case message, ok := <-c.sendChan:
			if !ok {
				return
			}
			c.ws.WriteMessage(websocket.BinaryMessage, message)
		case <-ticker.C:
			c.ws.WriteMessage(websocket.PingMessage, nil)
		}
	}
}

func (g *Gateway) readPump(c *Connection) {
	defer g.connections.Delete(c.playerId)

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			g.log.Log(log.LevelError, err)
			return
		}
		g.log.Log(log.LevelInfo, message)

		c.sendChan <- []byte("hello world")
	}
}
