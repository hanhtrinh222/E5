# Go WebSocket Chat with Redis
Ứng dụng này là server chat sử dụng Go và Redis để broadcast tin nhắn giữa các client thông qua WebSocket.

## Tính năng

- Giao tiếp real-time giữa nhiều client qua WebSocket.
- Broadcast tin nhắn bằng Redis Pub/Sub.
- Dễ dàng mở rộng và triển khai.

## Yêu cầu

- Go (>= 1.18)
- Redis server (Windows: dùng bản [không chính thức](https://github.com/tporadowski/redis/releases))
- (Tùy chọn) Docker nếu muốn chạy Redis nhanh

## Cài đặt Redis trên Windows

1. Tải file zip Redis từ [github.com/tporadowski/redis/releases](https://github.com/tporadowski/redis/releases)
2. Giải nén vào thư mục, ví dụ: `C:\Redis`
3. Mở cmd/PowerShell tại thư mục đó, chạy:  

   .\redis-server.exe
   hoặc dùng Docker:
   docker run -p 6379:6379 --name redis -d redis
  

## Chạy server Go

1. Clone repo về máy.
2. Cài đặt package cần thiết (nếu có):
   go mod tidy

3. Chạy server:
   go run main.go
   Mặc định server lắng nghe tại `ws://localhost:8080/ws`

## Kết nối client

Bạn có thể dùng file HTML demo sau:

```html
<!DOCTYPE html>
<html>
<body>
  <input id="msg" type="text" /><button onclick="send()">Gửi</button>
  <div id="chat"></div>
  <script>
    var ws = new WebSocket("ws://localhost:8080/ws");
    ws.onmessage = function(event) {
      document.getElementById("chat").innerHTML += `<div>${event.data}</div>`;
    };
    function send() {
      ws.send(document.getElementById("msg").value);
    }
  </script>
</body>
</html>
Mở nhiều tab để test broadcast.

- **Không đóng cửa sổ cmd** đang chạy Redis hoặc Go server, nếu không server sẽ tắt.
- Nếu cổng 8080/6379 bị chiếm, hãy kiểm tra và giải phóng bằng lệnh:
  netstat -ano | findstr :8080
  taskkill /PID <PID> /F


