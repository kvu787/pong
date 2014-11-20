package main

import (
	"math"
	"runtime"
	"time"
)

// data

func main() {

	/*****************/
	/*** LOAD DATA ***/
	/*****************/

	WINDOW = Rectangle_s{
		Position: Vector_s{X: 320, Y: 240},
		Width:    640,
		Height:   480,
	}
	SF_WINDOW = SetupWindow(int(WINDOW.Width), int(WINDOW.Height))
	FPS = 63
	FRAME_DURATION = SecondsToDuration(1.0 / float64(FPS))
	CURRENT_FRAME_START_TIME = time.Now()
	PLAYER =
		Paddle_s{
			Rectangle: Rectangle_s{
				Position: Vector_s{
					X: PADDLE_OFFSET,
					Y: WINDOW.Height / 2},
				Width:  25,
				Height: 100},
			Speed: 300}
	AI =
		Paddle_s{
			Rectangle: Rectangle_s{
				Position: Vector_s{
					X: WINDOW.Width - PADDLE_OFFSET,
					Y: WINDOW.Height / 2},
				Width:  25,
				Height: 100},
			Speed: 150}
	BALL =
		Ball_s{
			Circle: Circle_s{
				Position: Vector_s{
					X: WINDOW.Width / 2,
					Y: WINDOW.Height / 2},
				Radius: 10},
			Velocity: NewPolar(400, math.Pi*0.7)}
	INPUT =
		Input_s{
			IsUpArrowClicked:   false,
			IsDownArrowClicked: false,
			IsWindowClosed:     false}
	PADDLE_OFFSET = 50

	runtime.LockOSThread()

	/*********************/
	/*** RUN RELATIONS ***/
	/*********************/

	for {
		// handle timing
		CURRENT_FRAME_START_TIME = UpdateCurrentFrameStartTime()

		// handle input
		INPUT = UpdateInput(SF_WINDOW, INPUT)
		HandleWindowClose(INPUT)

		// update game objects
		PLAYER = MovePlayer(INPUT, FRAME_DURATION, PLAYER)
		AI = MoveAi(AI, BALL, FRAME_DURATION)
		BALL = CollidePlayerBall(PLAYER, BALL, 2)
		BALL = CollideBoundaryBall(WINDOW, BALL)
		BALL = ApplyVelocityBall(BALL, FRAME_DURATION)

		// rendering
		ClearWindow(Color_s{0, 0, 0, 255}, SF_WINDOW)
		RenderRectangle(PLAYER.Rectangle, SF_WINDOW)
		RenderRectangle(AI.Rectangle, SF_WINDOW)
		RenderCircle(BALL.Circle, SF_WINDOW)
		DisplayWindow(SF_WINDOW)

		// handle timing
		CURRENT_FRAME_DURATION = UpdateCurrentFrameDuration(CURRENT_FRAME_START_TIME)
		Sync(FRAME_DURATION, CURRENT_FRAME_DURATION)
	}
}
