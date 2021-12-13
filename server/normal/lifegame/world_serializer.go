package lifegame

import (
	"encoding/json"
	"fmt"
)

func JsonSerialize(world LifeGameWorld) []byte {
	arr := make([][]string, len(world))
	for i := range arr {
		arr[i] = make([]string, len(world[i]))
	}

	for i := range world {
		for j := range world[i] {
			arr[i][j] = world[i][j].String()
		}
	}

	b, _ := json.Marshal(arr)
	return b
}

//スライスが矩形でない場合はエラーを返す
func getWidthHeight(arr [][]string) (width, height int, err error) {
	colSet := map[int]struct{}{}
	for _, row := range arr {
		colSet[len(row)] = struct{}{}
	}

	if len(colSet) > 1 {
		return 0, 0, fmt.Errorf("矩形である必要があります")
	}

	return len(arr[0]), len(arr), nil
}

func JsonUnmarshal(b []byte) (LifeGameWorld, error) {
	var rawWorld [][]string
	if err := json.Unmarshal(b, &rawWorld); err != nil {
		return nil, err
	}

	width, height, err := getWidthHeight(rawWorld)
	if err != nil {
		return nil, err
	}

	generator := func(x, y int) ILifeState {
		switch rawWorld[y][x] {
		case AliveState{}.String():
			return AliveState{}
		case DeathState{}.String():
			return DeathState{}
		default:
			return AliveState{}
		}
	}

	return InitLifeGameWorld(width, height, generator), nil
}
