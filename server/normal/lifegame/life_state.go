package lifegame

import (
	"fmt"
	"math/rand"
)

type ILifeState interface {
	nextState(aroundState) ILifeState
	fmt.Stringer
}

type aroundState [8]*ILifeState

type AliveState struct {
}

func (s AliveState) String() string {
	return "alive"
}

func (s AliveState) nextState(around aroundState) ILifeState {
	aroundAliveNum := countAroundLifeState(around, AliveState{})
	if aroundAliveNum == 2 || aroundAliveNum == 3 {
		return AliveState{}
	}

	return DeathState{}
}

type DeathState struct{}

func (s DeathState) String() string {
	return "death"
}

func (s DeathState) nextState(around aroundState) ILifeState {
	aroundAliveNum := countAroundLifeState(around, AliveState{})
	if aroundAliveNum == 3 {
		return AliveState{}
	}

	return DeathState{}
}

func generateRandomState(x, y int) ILifeState {
	allLifeStateKind := []ILifeState{DeathState{}, AliveState{}}
	i := rand.Int() % len(allLifeStateKind)

	return allLifeStateKind[i]
}

func countAroundLifeState(aroundState aroundState, countValue ILifeState) int {
	count := 0
	for _, s := range aroundState {
		if s == nil {
			continue
		}

		if countValue == *s {
			count++
		}
	}

	return count
}
