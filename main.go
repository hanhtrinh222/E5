package main

import (
    "context"
    "log"
    "net/http"
    "sync"

    "github.com/gorilla/websocket"
    "github.com/go-redis/redis/v8"
)

var (
    upgrader    = websocket.Upgrader{CheckOrigin: func(r *http.Request) bool { return true }}
    clients     = make(map[*websocket.Conn]bool)
    mu          sync.Mutex
    redisClient *redis.Client
    ctx         = context.Background()
    channelName = "broadcast"
)

func main() {
    // Khởi tạo Redis client
    redisClient = redis.NewClient(&redis.Options{
        Addr: "localhost:6379", // Đổi lại nếu Redis ở server khác
    })

    // Lắng nghe Redis Pub/Sub, broadcast cho client
    go subscribeRedis()

    http.HandleFunc("/ws", handleWS)
    log.Println("WebSocket server started at :8080/ws")
    log.Fatal(http.ListenAndServe(":8080", nil))
}

// Xử lý kết nối WebSocket
func handleWS(w http.ResponseWriter, r *http.Request) {
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return
    }
    mu.Lock()
    clients[conn] = true
    mu.Unlock()

    defer func() {
        mu.Lock()
        delete(clients, conn)
        mu.Unlock()
        conn.Close()
    }()

    for {
        _, msg, err := conn.ReadMessage()
        if err != nil {
            break
        }
        // Publish lên Redis để broadcast cho toàn bộ instance
        err = redisClient.Publish(ctx, channelName, msg).Err()
        if err != nil {
            log.Println("Publish error:", err)
        }
    }
}

// Nhận message từ Redis và broadcast cho client
func subscribeRedis() {
    sub := redisClient.Subscribe(ctx, channelName)
    ch := sub.Channel()
    for msg := range ch {
        broadcast([]byte(msg.Payload))
    }
}

// Gửi message tới tất cả client đang kết nối
func broadcast(message []byte) {
    mu.Lock()
    defer mu.Unlock()
    for client := range clients {
        if err := client.WriteMessage(websocket.TextMessage, message); err != nil {
            client.Close()
            delete(clients, client)
        }
    }
}