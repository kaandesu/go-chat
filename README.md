# Go Chatroom Server

This is a simple chatroom server written in Go. It allows multiple users to connect, chat with each other, and join different chat rooms.

### Features

- **Multiple Users**: Connect multiple users to the server simultaneously.
- **Chat Rooms**: Users can join different chat rooms and communicate with each other.
- **Simple and Lightweight**: Written in Go, making it lightweight and easy to deploy.

### Usage

1. **Clone the Repository**:

   ```bash
   git clone https://github.com/kaandesu/go-chat.git
   cd go-chat
   ```

2. **Run the server**:

```bash
make run
```

3. **Connect to the Server**:
   Use a TCP client (e.g., telnet or netcat) to connect to the server:

```bash
nc localhost 3000
```

### Contributions

Contributions are welcome! Feel free to open issues or pull requests for any improvements or bug fixes.

---

_I will update this readme, once it reaches a stable point with its initial limited features and I refactor the file structure._
