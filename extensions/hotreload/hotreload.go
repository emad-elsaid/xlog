package hotreload

import (
	"fmt"
	"html/template"
	"log/slog"
	"sync"

	. "github.com/emad-elsaid/xlog"
	"github.com/gorilla/websocket"
)

var (
	upgrader     = websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024}
	clients      = make(map[*websocket.Conn]bool)
	clientsMutex sync.Mutex
)

func init() {
	RegisterExtension(Hotreload{})
}

type Hotreload struct{}

func (Hotreload) Name() string { return "hotreload" }
func (Hotreload) Init() {
	if Config.Readonly {
		return
	}

	Listen(Changed, NotifyPageChange)
	Get(`/+/hotreload`, handleWebSocket)
	RegisterWidget(AFTER_VIEW_WIDGET, 0, clientWidget)
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
			BadRequest(err.Error())(w, r)
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

const clientScript = `
    <script>
    (() => {
        const socketUrl = 'ws://'+window.location.host+'/+/hotreload';
        let socket = new WebSocket(socketUrl);
        socket.addEventListener('message', (evt) => {
            let data = JSON.parse(evt.data)
  			sessionStorage.setItem('scrollPosition', window.scrollY);
            window.location.href = data.url;
        });
    })();
    window.addEventListener('load', function() {
    	const scrollPosition = sessionStorage.getItem('scrollPosition');
    	if (scrollPosition !== null) {
    		window.scrollTo(0, parseInt(scrollPosition, 10));
    	}
    });
    </script>
    `

func clientWidget(p Page) template.HTML {
	return template.HTML(clientScript)
}
