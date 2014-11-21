package main

import (
	"fmt"
)

func main() {
	fmt.Println("Let's make Pong!")

	WINDOW_X = 640
	WINDOW_Y = 480
	SF_WINDOW = SetupWindow(WINDOW_X, WINDOW_Y)
	PLAYER = {Position: {20, 40}, Speed: 10, Dimensions: {100, 20}}

	for {
		// read window input
		INPUT = UpdateInput(SF_WINDOW)

		// relationship: INPUT -> "program"
		HandleWindowClose(INPUT)

		// relationship: INPUT -> update PLAYER
		PLAYER = UpdatePlayer(INPUT, PLAYER)

		// render window
		ClearWindow(SF_WINDOW)
		RenderPlayer(PLAYER, SF_WINDOW)
		DisplayWindow(SF_WINDOW)
	}
}
