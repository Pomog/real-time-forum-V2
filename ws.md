# WebSocket Utility Documentation
The WebSocket utility module is responsible for establishing and managing WebSocket connections within the application. It provides methods to connect, send messages, and disconnect from the WebSocket server. The module also handles incoming messages and triggers appropriate actions based on the message type.

## Features
  - **Connection Management:**
Establishes and maintains a WebSocket connection.
Automatically reconnects if the connection is lost or closed.
Ensures that a valid WebSocket connection is always used when sending messages.

  - **Message Handling:**
Listens for various types of messages from the server and triggers corresponding actions.
Handles token expiration by refreshing the token and re-authenticating.

  - **Error Handling:**
Displays errors when the connection fails or when invalid tokens are detected.
Alerts the user if their browser does not support WebSockets.

## Module Structure
  - **getConnection()**
The getConnection() function is responsible for creating or retrieving an active WebSocket connection. It returns a promise that resolves with the WebSocket connection instance.
Example:
```javascript
const connection = await getConnection();
```

## WebSocket Event Handlers:
  - **onopen:** Sends the authentication token to the server when the connection is established.
  - **onerror:** Displays an error message using Utils.showError() if the connection fails.
  - **onmessage:** Handles incoming messages from the server and triggers corresponding actions based on the message type.

### Example Event Handling:
```javascript
conn.onmessage = async function (evt) {
    let obj = JSON.parse(evt.data);
    switch (obj.type) {
        case "message":
            Chats.drawNewMessage(obj.body);
            break;
        case "messagesResponse":
            Chats.prependMessages(obj.body);
            break;
        // Additional cases for handling other message types
    }
};
```

### Ws Object
The Ws object provides methods for managing the WebSocket connection.

  - **connect()**
Establishes a WebSocket connection by calling getConnection() and storing the connection instance.
```javascript
Ws.connect();
```
  - **send(e)**
Sends a message over the WebSocket connection. Ensures that a connection is established before sending.
```javascript
Ws.send(JSON.stringify({ type: "message", body: "Hello, World!" }));
```

  - **disconnect()**
Closes the WebSocket connection.**
```javascript
Ws.disconnect();
```

## External Dependencies
- [Chats](chats.md): Manages the chat interface, including rendering messages, chats, and online users.
- [Utils](utils.md): Provides utility functions, including error handling.
- [Fetcher](fetcher.md): Handles HTTP requests, including token refresh logic.
- [Router](router.md): Manages routing within the application.

## Usage Example
```javascript
import Ws from './path/to/Ws.js';

// Establish a connection
await Ws.connect();

// Send a message
await Ws.send(JSON.stringify({ type: "message", body: "Hello, World!" }));

// Disconnect
await Ws.disconnect();
```
