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
		Position: Vector_s{X: 400, Y: 300},
		Width:    800,
		Height:   600,
	}
	SF_WINDOW = SetupWindow(int(WINDOW.Width), int(WINDOW.Height))
	FPS = 120
	FRAME_DURATION = SecondsToDuration(1.0 / float64(FPS))
	CURRENT_FRAME_START_TIME = time.Now()
	PADDLE_OFFSET = 60
	PLAYER =
		Paddle_s{
			Rectangle: Rectangle_s{
				Position: Vector_s{
					X: PADDLE_OFFSET,
					Y: WINDOW.Height / 2},
				Width:  10,
				Height: 80},
			Speed: 500}
	AI =
		Paddle_s{
			Rectangle: Rectangle_s{
				Position: Vector_s{
					X: WINDOW.Width - PADDLE_OFFSET,
					Y: WINDOW.Height / 2},
				Width:  10,
				Height: 80},
			Speed: 300}
	BALL =
		Ball_s{
			Circle: Circle_s{
				Position: Vector_s{
					X: WINDOW.Width / 2,
					Y: WINDOW.Height / 2},
				Radius: 5},
			Velocity: NewPolar(700, math.Pi*0.05)}
	INPUT =
		Input_s{
			IsUpArrowClicked:   false,
			IsDownArrowClicked: false,
			IsWindowClosed:     false}

	runtime.LockOSThread()

	/*********************/
	/*** RUN RELATIONS ***/
	/*********************/

	for { // game loop
		// handle timing
		CURRENT_FRAME_START_TIME = UpdateCurrentFrameStartTime()

		// handle input
		INPUT = UpdateInput(SF_WINDOW, INPUT)
		HandleWindowClose(INPUT)

		// update game objects
		PLAYER = MovePlayer(INPUT, FRAME_DURATION, PLAYER)
		PLAYER = KeepPlayerInBoundary(PLAYER, WINDOW)
		AI = MoveAi(AI, BALL, FRAME_DURATION)
		BALL = CollidePaddleBall(PLAYER, BALL, 4, DegreesToRadians(20.0))
		BALL = CollidePaddleBall(AI, BALL, 4, DegreesToRadians(20.0))
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
		Sleep(FRAME_DURATION, CURRENT_FRAME_DURATION)
	}
}

/* Pong collision model 3

if paddle is ball_diameter away from the boundary
after collisions, the ball must be CLEAR of the danger zone

*/
