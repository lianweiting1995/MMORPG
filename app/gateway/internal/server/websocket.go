package server

import (
	pb "MMORPG/api/websocket/v1"
	"MMORPG/app/gateway/internal/service"
	"context"
	"encoding/json"
	"fmt"
	"net/http"
	"sync"
	"sync/atomic"
	"time"

	"github.com/go-kratos/kratos/v2/log"
	"github.com/google/uuid"
	"github.com/gorilla/websocket"
	"google.golang.org/protobuf/proto"
)

type Connection struct {
	ws         *websocket.Conn
	sendChan   chan []byte
	playerId   string
	zoneID     int
	lastActive int64 // 记录最后的活跃时间
}

type Gateway struct {
	connections *sync.Map             // playerID -> *Connection
	zoneBuckets [256]*sync.Map        // 按区域分桶
	upgrader    websocket.Upgrader    // 连接标记
	log         log.Logger            // 日志记录
	loginSrv    *service.LoginService // 登录服务
}

func InitGateway(log log.Logger, loginSrv *service.LoginService) *Gateway {
	g := &Gateway{
		upgrader: websocket.Upgrader{
			ReadBufferSize:    1024,
			WriteBufferSize:   1024,
			EnableCompression: true,
		},
		log:         log,
		connections: &sync.Map{},
		loginSrv:    loginSrv,
	}
	// 初始化所有的桶
	for i := range g.zoneBuckets {
		g.zoneBuckets[i] = &sync.Map{}
	}

	return g
}

// calculateZone 计算得出区域位置
func calculateZone(playerID string) int {
	sum := 0
	for _, b := range []byte(playerID) {
		sum += int(b)
	}

	return sum % 256
}

func (g *Gateway) HandleConnection(ctx context.Context, w http.ResponseWriter, r *http.Request) {
	conn, err := g.upgrader.Upgrade(w, r, nil)
	if err != nil {
		g.log.Log(log.LevelError, "upgrade error:", err)
		return
	}
	conn.SetReadLimit(64 * 1024)                           // 64KB 设置读取最大的消息长度
	conn.SetReadDeadline(time.Now().Add(60 * time.Second)) // 设置读取超时时间
	conn.SetPongHandler(func(string) error {
		conn.SetReadDeadline(time.Now().Add(10 * time.Second))
		return nil
	})

	c := &Connection{
		ws:       conn,
		sendChan: make(chan []byte, 100),
		playerId: uuid.New().String(),
	}
	g.connections.Store(c.playerId, c)
	zone := calculateZone(c.playerId)
	g.zoneBuckets[zone].Store(c.playerId, c)

	go g.readPump(ctx, c)
	go g.writePump(ctx, c)
}

func (g *Gateway) writePump(ctx context.Context, c *Connection) {
	ticker := time.NewTicker(15 * time.Second)
	defer func() {
		ticker.Stop()
		close(c.sendChan)
		c.ws.Close()
	}()

	for {
		select {
		case message, ok := <-c.sendChan:
			if !ok {
				c.ws.WriteMessage(websocket.CloseMessage, nil)
				return
			}
			// 解析提交的数据
			var msg pb.GatewayMessage
			if err := proto.Unmarshal(message, &msg); err != nil {
				g.log.Log(log.LevelError, "%s", err)
				return
			}
			var re string
			switch msg.Type {
			case pb.MsgType_LOGIN:
				var input pb.LoginInput
				err := json.Unmarshal([]byte(msg.Data), &input)
				if err != nil {
					re = fmt.Sprintf("%v", err)
				} else {
					output, err := g.loginSrv.LoginWithPhone(ctx, &input)
					if err != nil {
						re = fmt.Sprintf("%v", err)
					} else {
						re = fmt.Sprintf("%+v", output)
					}
				}
			default:
				re = "invalid data"
			}
			// 设置读写超时保护
			c.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.ws.WriteMessage(websocket.BinaryMessage, []byte(re)); err != nil {
				g.log.Log(log.LevelError, "%s", err)
				return
			}
		case <-ticker.C:
			c.ws.SetWriteDeadline(time.Now().Add(10 * time.Second))
			if err := c.ws.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

func (g *Gateway) readPump(_ context.Context, c *Connection) {
	defer func() {
		g.connections.Delete(c.playerId)
		zone := calculateZone(c.playerId)
		g.zoneBuckets[zone].Delete(c.playerId)
		c.ws.Close()
	}()

	for {
		_, message, err := c.ws.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseNormalClosure) {
				g.log.Log(log.LevelError, err)
			}
			return
		}

		c.sendChan <- []byte(message)
	}
}

// checkInactiveConnections 定期检查连接
func (g *Gateway) checkInactiveConnections() {
	ticker := time.NewTicker(time.Minute)
	defer ticker.Stop()

	for range ticker.C {
		now := time.Now().Unix()
		g.connections.Range(func(key, value any) bool {
			c := value.(*Connection)
			if atomic.LoadInt64(&c.lastActive) < now-300 {
				c.ws.Close()
			}

			return true
		})
	}
}
