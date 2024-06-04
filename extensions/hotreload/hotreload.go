package hotreload

import (
    "fmt"
    "html/template"
    "log"
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
    Listen(Changed, NotifyPageChange)
    Get(`/ws`, handleWebSocket)
    RegisterWidget(AFTER_VIEW_WIDGET, 0, clientWidget)
}


func NotifyPageChange(p Page) error {
    if !p.Exists() {
        return nil
    }

    message := map[string]string{"url": fmt.Sprintf("/%s", p.Name())}
    log.Printf("Page %s changed from notifiy", p.Name())

    clientsMutex.Lock()
    defer clientsMutex.Unlock()

    for client := range clients {
        err := client.WriteJSON(message)
        if err != nil {
            log.Printf("Error sending message to client: %v\n", err)
            client.Close()
            delete(clients, client)
        }
    }
    return nil
}

func nop(w Response, r Request) {
}
func handleWebSocket(w Response, r Request) Output {
    if !READONLY {
        return NoContent()
    }
    conn, err := upgrader.Upgrade(w, r, nil)
    if err != nil {
        log.Println("Upgrade error:", err)
        return BadRequest(err.Error())
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
            if err != nil || mt ==  websocket.CloseMessage {
                break
            }
        }
    }()

    clientsMutex.Lock()
    clients[conn] = true
    clientsMutex.Unlock()
    log.Printf("Connection arrived. curr cons: %d", len(clients))
    return nop
}

// TODO use same HOST and PORT than server
const clientScript = `
    <script>
    (() => {
        const socketUrl = 'ws://'+window.location.host+'/ws';
        let socket = new WebSocket(socketUrl);
        socket.addEventListener('message', (evt) => {
            let data = JSON.parse(evt.data)
            window.location.href = data.url;
        });
    })();
    </script>
    `

func clientWidget(p Page) template.HTML {
    // return template.HTML(fmt.Sprint(clientScript, template.JSEscapeString(p.Name())))
    if !READONLY {
        return template.HTML(clientScript)
    }

    return ""
}
