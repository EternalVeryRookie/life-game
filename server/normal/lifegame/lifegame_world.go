package lifegame

import (
	"math/rand"
	"sync"
)

type Point2D struct {
	x, y int
}

func (p Point2D) Add(dx, dy int) Point2D {
	return Point2D{
		x: p.x + dx,
		y: p.y + dy,
	}
}

type LifeGameWorld [][]ILifeState

func GenerateRandomInitWorld(width, height int, seed int64) LifeGameWorld {
	rand.Seed(seed)
	return InitLifeGameWorld(width, height, generateRandomState)
}

func InitEmptyLifeGameWorld(width, height int) LifeGameWorld {
	world := make(LifeGameWorld, height)
	for i := range world {
		world[i] = make([]ILifeState, width)
	}

	return world
}

func InitLifeGameWorld(width, height int, generator func(x, y int) ILifeState) LifeGameWorld {
	world := make(LifeGameWorld, height)
	for y := range world {
		world[y] = make([]ILifeState, width)
		for x := range world[y] {
			world[y][x] = generator(x, y)
		}
	}

	return world
}

func getLiveState(world LifeGameWorld, p Point2D) *ILifeState {
	y := p.y
	if y < 0 {
		y = len(world) - 1
	} else if y >= len(world) {
		y = 0
	}

	x := p.x
	if x < 0 {
		x = len(world[y]) - 1
	} else if x >= len(world[y]) {
		x = 0
	}

	return &world[y][x]
}

func getAroundLiveState(world LifeGameWorld, p Point2D) aroundState {
	return aroundState{
		getLiveState(world, p.Add(0, 1)),
		getLiveState(world, p.Add(0, -1)),
		getLiveState(world, p.Add(-1, 0)),
		getLiveState(world, p.Add(1, 0)),
		getLiveState(world, p.Add(1, 1)),
		getLiveState(world, p.Add(1, -1)),
		getLiveState(world, p.Add(-1, 1)),
		getLiveState(world, p.Add(-1, -1)),
	}
}

func NextStepWorld(currentWorld LifeGameWorld) LifeGameWorld {
	nextWorld := InitEmptyLifeGameWorld(len(currentWorld[0]), len(currentWorld))
	wait := sync.WaitGroup{}
	for y := range currentWorld {
		for x := range currentWorld[y] {
			wait.Add(1)
			go func(y, x int) {
				defer wait.Done()
				around := getAroundLiveState(currentWorld, Point2D{x: x, y: y})
				nextWorld[y][x] = currentWorld[y][x].nextState(around)
			}(y, x)
		}
	}

	wait.Wait()

	return nextWorld
}
