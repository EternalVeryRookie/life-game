package lifegame

func send(c chan LifeGameWorld, world LifeGameWorld) (err error) {
	defer func() {
		err, _ = recover().(error)
	}()

	c <- world
	return
}

func SimulateAsync(init LifeGameWorld) chan LifeGameWorld {
	const bufferSize = 10
	c := make(chan LifeGameWorld, bufferSize)

	go func() {
		current := init
		var err error = nil
		for err == nil {
			current = NextStepWorld(current)
			err = send(c, current)
		}
	}()

	return c
}

func sendByte(c chan []byte, b []byte) (err error) {
	defer func() {
		err, _ = recover().(error)
	}()

	c <- b
	return
}

func SimulateResultSerialize(init LifeGameWorld) chan []byte {
	const bufferSize = 10
	c := make(chan []byte, bufferSize)

	go func() {
		current := init
		var err error = nil
		for err == nil {
			current = NextStepWorld(current)
			b := JsonSerialize(current)
			err = sendByte(c, b)
		}
	}()

	return c
}
