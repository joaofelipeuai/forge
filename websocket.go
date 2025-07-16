package main

import (
	"crypto/sha1"
	"encoding/base64"
	"fmt"
	"net/http"
	"strings"
)

// WebSocket constants
const (
	websocketMagicString = "258EAFA5-E914-47DA-95CA-C5AB0DC85B11"
)

// WebSocketHandler represents a WebSocket handler function
type WebSocketHandler func(*WebSocketConnection)

// WebSocketConnection represents a WebSocket connection
type WebSocketConnection struct {
	conn   http.ResponseWriter
	req    *http.Request
	closed bool
}

// WebSocket upgrade middleware
func (f *Forge) WebSocket(pattern string, handler WebSocketHandler) {
	f.GET(pattern, func(c *Context) error {
		return f.upgradeWebSocket(c, handler)
	})
}

// upgradeWebSocket handles the WebSocket upgrade process
func (f *Forge) upgradeWebSocket(c *Context, handler WebSocketHandler) error {
	// Check if it's a WebSocket upgrade request
	if !isWebSocketUpgrade(c.Request) {
		return c.String(400, "Bad Request: Not a WebSocket upgrade")
	}

	// Get the WebSocket key
	key := c.Request.Header.Get("Sec-WebSocket-Key")
	if key == "" {
		return c.String(400, "Bad Request: Missing Sec-WebSocket-Key")
	}

	// Generate accept key
	acceptKey := generateAcceptKey(key)

	// Set WebSocket headers
	c.Response.Header().Set("Upgrade", "websocket")
	c.Response.Header().Set("Connection", "Upgrade")
	c.Response.Header().Set("Sec-WebSocket-Accept", acceptKey)
	c.Response.WriteHeader(101) // Switching Protocols

	// Create WebSocket connection
	wsConn := &WebSocketConnection{
		conn: c.Response,
		req:  c.Request,
	}

	// Handle the WebSocket connection
	handler(wsConn)

	return nil
}

// isWebSocketUpgrade checks if the request is a WebSocket upgrade
func isWebSocketUpgrade(r *http.Request) bool {
	return strings.ToLower(r.Header.Get("Upgrade")) == "websocket" &&
		strings.Contains(strings.ToLower(r.Header.Get("Connection")), "upgrade")
}

// generateAcceptKey generates the WebSocket accept key
func generateAcceptKey(key string) string {
	h := sha1.New()
	h.Write([]byte(key + websocketMagicString))
	return base64.StdEncoding.EncodeToString(h.Sum(nil))
}

// Send sends a text message to the WebSocket client
func (ws *WebSocketConnection) Send(message string) error {
	if ws.closed {
		return fmt.Errorf("connection is closed")
	}

	// Simple text frame (this is a basic implementation)
	frame := []byte{0x81} // FIN=1, opcode=1 (text)
	
	msgLen := len(message)
	if msgLen < 126 {
		frame = append(frame, byte(msgLen))
	} else if msgLen < 65536 {
		frame = append(frame, 126)
		frame = append(frame, byte(msgLen>>8), byte(msgLen&0xFF))
	} else {
		frame = append(frame, 127)
		for i := 7; i >= 0; i-- {
			frame = append(frame, byte(msgLen>>(i*8)))
		}
	}
	
	frame = append(frame, []byte(message)...)
	
	_, err := ws.conn.Write(frame)
	return err
}

// Close closes the WebSocket connection
func (ws *WebSocketConnection) Close() error {
	if ws.closed {
		return nil
	}
	
	ws.closed = true
	// Send close frame
	closeFrame := []byte{0x88, 0x00} // FIN=1, opcode=8 (close)
	_, err := ws.conn.Write(closeFrame)
	return err
}

// WebSocket middleware for broadcasting
func WebSocketBroadcast() *WebSocketBroadcaster {
	return &WebSocketBroadcaster{
		connections: make(map[*WebSocketConnection]bool),
	}
}

type WebSocketBroadcaster struct {
	connections map[*WebSocketConnection]bool
}

func (wb *WebSocketBroadcaster) AddConnection(conn *WebSocketConnection) {
	wb.connections[conn] = true
}

func (wb *WebSocketBroadcaster) RemoveConnection(conn *WebSocketConnection) {
	delete(wb.connections, conn)
}

func (wb *WebSocketBroadcaster) Broadcast(message string) {
	for conn := range wb.connections {
		if !conn.closed {
			conn.Send(message)
		}
	}
}