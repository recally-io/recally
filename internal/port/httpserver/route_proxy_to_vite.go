package httpserver

import (
	"net/http"
	"net/http/httputil"
	"net/url"
	"strings"
	"vibrain/internal/pkg/logger"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

var upgrader = websocket.Upgrader{
	CheckOrigin: func(r *http.Request) bool {
		return true // TODO: implement origin checks to improve security
	},
}

// reactReverseProxy is a reverse proxy for vite server
func reactReverseProxy(c echo.Context) error {
	remote, _ := url.Parse("http://localhost:5173")
	proxy := httputil.NewSingleHostReverseProxy(remote)
	proxy.Director = func(req *http.Request) {
		req.Header = c.Request().Header
		req.Host = remote.Host
		req.URL = c.Request().URL
		req.URL.Scheme = remote.Scheme
		req.URL.Host = remote.Host
	}

	if isWebSocketRequest(c.Request()) {
		websocketProxyHandler(c.Response().Writer, c.Request())
	} else {
		proxy.ServeHTTP(c.Response().Writer, c.Request())
	}
	return nil
}

func websocketProxyHandler(w http.ResponseWriter, r *http.Request) {
	ctx := r.Context()
	// Parse the target WebSocket URL
	targetUrl := "ws://localhost:5173"
	targetURL, err := url.Parse(targetUrl)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to parse the target WebSocket URL", "err", err)
		return
	}
	proxyHeader := http.Header{}
	for k, v := range r.Header {
		proxyHeader[k] = v
	}
	// Dial the target WebSocket server
	targetConn, _, err := websocket.DefaultDialer.Dial(targetURL.String()+r.URL.Path, nil)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to connect to target WebSocket server", "err", err)
		return
	}
	defer targetConn.Close()

	// Upgrade the incoming connection to a WebSocket connection
	clientConn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		logger.FromContext(ctx).Error("Failed to upgrade the connection", "err", err)
		return
	}
	defer clientConn.Close()

	// Channel to signal that one of the connections has closed
	done := make(chan struct{})

	// Forward messages from client to target server
	go func() {
		defer close(done)
		for {
			messageType, message, err := clientConn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
					logger.FromContext(ctx).Info("Client closed connection")
				} else {
					logger.FromContext(ctx).Error("Failed to read message from client", "err", err)
				}
				return
			}
			err = targetConn.WriteMessage(messageType, message)
			if err != nil {
				logger.FromContext(ctx).Error("Failed to forward message to target", "err", err)
				return
			}
		}
	}()

	// Forward messages from target server to client
	go func() {
		defer close(done)
		for {
			messageType, message, err := targetConn.ReadMessage()
			if err != nil {
				if websocket.IsCloseError(err, websocket.CloseNormalClosure, websocket.CloseGoingAway, websocket.CloseNoStatusReceived) {
					logger.FromContext(ctx).Info("Target server closed connection")
				} else {
					logger.FromContext(ctx).Error("Failed to read message from target server", "err", err)
				}
				return
			}
			err = clientConn.WriteMessage(messageType, message)
			if err != nil {
				logger.FromContext(ctx).Error("Failed to forward message to client", "err", err)
				return
			}
		}
	}()

	// Wait for one of the connections to close
	<-done

	// Perform cleanup if needed
}

func isWebSocketRequest(r *http.Request) bool {
	connHeader := strings.ToLower(r.Header.Get("Connection"))
	upgradeHeader := strings.ToLower(r.Header.Get("Upgrade"))
	return strings.Contains(connHeader, "upgrade") &&
		strings.Contains(upgradeHeader, "websocket")
}
