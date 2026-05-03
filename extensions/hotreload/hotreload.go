package hotreload

import (
	"fmt"
	"html/template"
	"log/slog"
	"sync"

	_ "embed"

	"github.com/emad-elsaid/xlog"
	"github.com/gorilla/websocket"
)

var (
	upgrader     = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	clients      = make(map[*websocket.Conn]bool)
	clientsMutex sync.Mutex
)

func init() {
	xlog.RegisterExtension(Hotreload{})
}

type Hotreload struct{}

func (Hotreload) Name() string { return "hotreload" }
func (Hotreload) Init() {
	if xlog.Config.Readonly {
		return
	}

	xlog.Listen(xlog.PageChanged, NotifyPageChange)
	xlog.Get(`/+/hotreload`, handleWebSocket)
	xlog.RegisterWidget(xlog.WidgetAfterView, 0, clientWidget)
}

func NotifyPageChange(p xlog.Page) error {
	if !p.Exists() {
		return nil
	}

	message := map[string]string{"url": fmt.Sprintf("/%s", p.Name())}

	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range clients {
		err := client.WriteJSON(message)
		if err != nil {
			if err := client.Close(); err != nil {
				slog.Error("Failed to close websocket client", "error", err)
			}
			delete(clients, client)
		}
	}
	return nil
}

func handleWebSocket(r xlog.Request) xlog.Output {
	return func(w xlog.Response, r xlog.Request) {
		conn, err := upgrader.Upgrade(w, r, nil)
		if err != nil {
			slog.Error("Failed to upgrade", "error", err)
			xlog.BadRequest(err.Error())(w, r)
		}

		// keep connection open
		go func() {
			defer func() {
				clientsMutex.Lock()
				delete(clients, conn)
				clientsMutex.Unlock()
				if err := conn.Close(); err != nil {
					slog.Error("Failed to close websocket connection", "error", err)
				}
			}()

			for {
				mt, _, err := conn.ReadMessage()
				if err != nil || mt == websocket.CloseMessage {
					break
				}
			}
		}()

		clientsMutex.Lock()
		clients[conn] = true
		clientsMutex.Unlock()
	}
}

//go:embed script.html
var clientScript string

func clientWidget(p xlog.Page) template.HTML {
	return template.HTML(clientScript)
}
