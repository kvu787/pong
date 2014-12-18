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
	PLAYER_1 =
		Paddle_s{
			Rectangle: Rectangle_s{
				Position: Vector_s{
					X: PADDLE_OFFSET,
					Y: WINDOW.Height / 2},
				Width:  10,
				Height: 80},
			Speed: 500}
	PLAYER_2 =
		Paddle_s{
			Rectangle: Rectangle_s{
				Position: Vector_s{
					X: WINDOW.Width - PADDLE_OFFSET,
					Y: WINDOW.Height / 2},
				Width:  10,
				Height: 80},
			Speed: 500}
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
			IsWClicked:         false,
			IsSClicked:         false,
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
		PLAYER_1 = MovePlayer(INPUT, FRAME_DURATION, PLAYER_1, 1)
		PLAYER_1 = KeepPlayerInBoundary(PLAYER_1, WINDOW)
		PLAYER_2 = MovePlayer(INPUT, FRAME_DURATION, PLAYER_2, 2)
		PLAYER_2 = KeepPlayerInBoundary(PLAYER_2, WINDOW)
		BALL = CollidePaddleBall(PLAYER_1, BALL, 4, DegreesToRadians(20.0))
		BALL = CollidePaddleBall(PLAYER_2, BALL, 4, DegreesToRadians(20.0))
		BALL = CollideBoundaryBall(WINDOW, BALL)
		BALL = ApplyVelocityBall(BALL, FRAME_DURATION)

		// rendering
		ClearWindow(Color_s{0, 0, 0, 255}, SF_WINDOW)
		RenderRectangle(PLAYER_1.Rectangle, SF_WINDOW)
		RenderRectangle(PLAYER_2.Rectangle, SF_WINDOW)
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
