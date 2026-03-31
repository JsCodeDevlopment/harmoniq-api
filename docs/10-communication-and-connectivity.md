# 📨 Communication and Connectivity

The Gost Framework provides a robust infrastructure for asynchronous communication (RabbitMQ) and real-time communication (Websockets), as well as an external event dispatch system (Webhooks) with delivery guarantees.

## 1. RabbitMQ (Asynchronous Messaging)

RabbitMQ is used to decouple heavy processes from the main HTTP request, ensuring scalability and resilience.

### ⚙️ Configuration (`src/config/rabbitmq.go`)
This file manages the persistent connection with the broker. It uses environment variables to connect and provides a global `RabbitMQConn` object.
- **Initialization**: Automatically called during the app setup.
- **Resilience**: Logs errors if the connection fails, allowing the app to start without RabbitMQ if the functionality is optional.

### 📤 Producer (`src/common/messaging/producer.go`)
A utility for sending messages to any queue or exchange.
- **`PublishMessage(exchange, routingKey, payload)`**: Automatically converts your `interface{}` (struct/map) to JSON and sends it to RabbitMQ with persistence enabled.

**Usage example:**
```go
err := messaging.PublishMessage("", "orders_queue", map[string]interface{}{
    "order_id": 123,
    "status": "processing",
})
```

### 📥 Consumer Scaffolding (`src/common/messaging/consumer.go`)
Basis for creating "Workers" that process messages in the background.
- **`RegisterConsumer(queue, handler)`**: Automatically declares the queue (durable) and starts a goroutine to listen for messages.
- **`SimpleAckHandler`**: A helper that handles Ack (confirmation) or Nack (failure) automatically based on the return of your business logic.

---

## 🌐 2. Websockets (Real Time)

Allows low-latency bidirectional communication between the server and clients.

### 🏗️ Central Hub (`src/modules/ws/hub.go`)
The "brain" of the connections. It manages who is connected and coordinates message delivery.
- **Registration/Unregistration**: Manages client entry and exit securely (thread-safe).
- **Broadcast**: `BroadcastJSON` method to send a message to **all** connected users simultaneously.

### 🎮 Controller (`src/modules/ws/ws.controller.go`)
Exposes the `/api/v1/ws` endpoint.
- **Upgrade**: Transforms a standard HTTP request into a persistent Websocket connection.
- **Pumps**: Starts read and write loops to keep the connection alive and process data.

**How to use it in the Frontend:**
```javascript
const socket = new WebSocket('ws://localhost:3000/api/v1/ws');
socket.onmessage = (event) => {
    console.log('Message from server:', JSON.parse(event.data));
};
```

---

## 🪝 3. Webhooks & Retries

Allows your system to notify other systems securely.

### 🚀 Webhook Dispatcher (`src/common/utils/webhook.go`)
- **HMAC Signature**: Every webhook sent contains an `X-Gost-Signature` header. This allows the receiver to verify that the message truly came from your server using a shared secret key.
- **`SendWebhook(url, secret, event, data)`**: Executes the POST request with the structured payload.

### 🔄 Retry Worker (`src/common/messaging/webhook_worker.go`)
Intelligent integration between Webhooks and RabbitMQ.
- **Resilience**: If the target server is down, Gost does not lose the event. It places it in a retry queue.
- **Exponential Backoff**: The system waits an increasing amount of time between attempts (1s, 2s, 3s...) up to a maximum of 5 attempts.

**Usage example:**
```go
// Dispatches and, if it fails, schedules an automatic retry via RabbitMQ
messaging.DispatchWebhookWithRetry(
    "https://client.com/callback", 
    "my_secret_key", 
    "user.created", 
    userData,
)
```

---

## 💡 Recommended Workflow

1. **Define an Event**: What happened? (e.g., Payment Approved).
2. **Send a Webhook**: Notify the client's system.
3. **Notify the Websocket**: Update the user interface in real-time without refreshing.
4. **Process in the Background**: If there is a heavy task (generating a PDF), send it to a queue via RabbitMQ.
