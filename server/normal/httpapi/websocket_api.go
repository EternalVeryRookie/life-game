package httpapi

import (
	"fmt"
	"net/http"
	"strings"

	"github.com/gorilla/websocket"
	"github.com/labstack/echo/v4"
)

type requestType string

const (
	start requestType = "start"
	next  requestType = "next"
	stop  requestType = "stop"
	end   requestType = "end"
	reset requestType = "reset"
)

func parseRequest(message string) (method requestType, jsonStr string, err error) {
	methodList := [5]requestType{start, next, stop, end, reset}
	for _, method := range methodList {
		if strings.HasPrefix(message, string(method)) {
			return method, strings.TrimLeft(message, string(method)), nil
		}
	}

	return "", "", fmt.Errorf("不明なリクエスト: %s", message)
}

func Simulate(c echo.Context) error {
	upgrader := websocket.Upgrader{ReadBufferSize: 1024, WriteBufferSize: 1024, CheckOrigin: func(r *http.Request) bool {
		return true
	}}

	conn, err := upgrader.Upgrade(c.Response(), c.Request(), nil)
	if err != nil {
		return c.JSON(http.StatusForbidden, map[string]interface{}{"error": err.Error()})
	}

	var state iState = &initState{}
	//stateを即時評価してはいけないためdeferでcloseは行わない
	defer conn.Close()

	for {
		_, b, err := conn.ReadMessage()
		if err != nil {
			break
		}

		request, body, err := parseRequest(string(b))
		if err != nil {
			break
		}

		state = state.recieve(request, body, conn)
	}

	state.close()
	return nil
}
