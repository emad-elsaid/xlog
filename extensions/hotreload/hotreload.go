package hotreload

import (
	"fmt"
	"html/template"
	"log/slog"
	"sync"

	_ "embed"

	"github.com/emad-elsaid/xlog"
	. "github.com/emad-elsaid/xlog"
	"github.com/gorilla/websocket"
)

var (
	upgrader     = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	clients      = make(map[*websocket.Conn]bool)
	clientsMutex sync.Mutex
)

func init() {
	app := GetApp()
	app.RegisterExtension(Hotreload{})
}

type Hotreload struct{}

func (Hotreload) Name() string { return "hotreload" }
func (Hotreload) Init(app *App) {
	if app.GetConfig().Readonly {
		return
	}

	app.Listen(PageChanged, NotifyPageChange)
	app.Get(`/+/hotreload`, handleWebSocket)
	app.RegisterWidget(WidgetAfterView, 0, clientWidget)
}

func NotifyPageChange(p Page) error {
	if !p.Exists() {
		return nil
	}

	message := map[string]string{"url": fmt.Sprintf("/%s", p.Name())}

	clientsMutex.Lock()
	defer clientsMutex.Unlock()

	for client := range clients {
		err := client.WriteJSON(message)
		if err != nil {
			client.Close()
			delete(clients, client)
		}
	}
	return nil
}

func handleWebSocket(r Request) Output {
	return func(w Response, r Request) {
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
				conn.Close()
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

func clientWidget(p Page) template.HTML {
	return template.HTML(clientScript)
}
