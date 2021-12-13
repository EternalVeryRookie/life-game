package httpapi

import (
	"life-game/lifegame"

	"github.com/gorilla/websocket"
)

type iState interface {
	recieve(requestType, string, *websocket.Conn) iState
	close()
}

type initState struct{}

func (s *initState) recieve(request requestType, body string, conn *websocket.Conn) iState {
	//TODO: サイズ制限
	if request == start {
		initWorld, err := lifegame.JsonUnmarshal([]byte(body))
		if err != nil {
			return s
		}

		q := lifegame.SimulateResultSerialize(initWorld)
		conn.WriteMessage(websocket.TextMessage, <-q)
		return newSimulationingState(initWorld, q)
	}

	return s
}

func (s *initState) close() {}

type simulationingState struct {
	initWorld  lifegame.LifeGameWorld
	stateQueue chan []byte
}

func newSimulationingState(init lifegame.LifeGameWorld, stateQueue chan []byte) *simulationingState {
	return &simulationingState{
		initWorld:  init,
		stateQueue: stateQueue,
	}
}

func (s *simulationingState) recieve(request requestType, body string, conn *websocket.Conn) iState {
	if request == stop {
		return &suspendState{s.initWorld, s.stateQueue}
	}

	if request == end {
		close(s.stateQueue)
		conn.WriteMessage(websocket.TextMessage, lifegame.JsonSerialize(s.initWorld))
		return &initState{}
	}

	if request == next {
		conn.WriteMessage(websocket.TextMessage, <-s.stateQueue)
		return s
	}

	return s
}

func (s *simulationingState) close() {
	defer func() {
		recover()
	}()

	close(s.stateQueue)
}

type suspendState struct {
	initWorld  lifegame.LifeGameWorld
	stateQueue chan []byte
}

func (s *suspendState) recieve(request requestType, body string, conn *websocket.Conn) iState {
	switch request {
	case start:
		conn.WriteMessage(websocket.TextMessage, <-s.stateQueue)
		return &simulationingState{initWorld: s.initWorld, stateQueue: s.stateQueue}
	case end:
		close(s.stateQueue)
		conn.WriteMessage(websocket.TextMessage, lifegame.JsonSerialize(s.initWorld))
		return &initState{}
	}

	return s
}

func (s *suspendState) close() {
	defer func() {
		recover()
	}()

	close(s.stateQueue)
}
