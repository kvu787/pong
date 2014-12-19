package main

import (
	"fmt"
	"runtime"
	"time"
)

// data

func main() {

	/*****************/
	/*** LOAD DATA ***/
	/*****************/

	BALL_START_VELOCITY = 350
	BALL_MAX_VELOCITY = 750
	BALL_ACCELERATION = 20
	WINDOW = Rectangle_s{
		Position: Vector_s{X: 400, Y: 300},
		Width:    800,
		Height:   600,
	}
	SF_WINDOW = SetupWindow(int(WINDOW.Width), int(WINDOW.Height))
	FPS = 65
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
			Rectangle: Rectangle_s{
				Position: Vector_s{
					X: WINDOW.Width / 2,
					Y: WINDOW.Height / 2},
				Width:  10,
				Height: 10},
			Velocity: NewPolar(BALL_START_VELOCITY, DegreesToRadians(173.5))}
	INPUT =
		Input_s{
			IsUpArrowClicked:   false,
			IsDownArrowClicked: false,
			IsWClicked:         false,
			IsSClicked:         false,
			IsWindowClosed:     false}
	PLAYER_1_SCORE = 0
	PLAYER_2_SCORE = 0
	HAS_SCORED = false
	HAS_HIT_PADDLE_1 = false
	HAS_HIT_PADDLE_2 = false

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
		BALL, HAS_HIT_PADDLE_1 = CollidePaddleBall(PLAYER_1, BALL)
		BALL, HAS_HIT_PADDLE_2 = CollidePaddleBall(PLAYER_2, BALL)
		if HAS_HIT_PADDLE_1 || HAS_HIT_PADDLE_2 {
			BALL = ClampBall(BALL, PLAYER_1, PLAYER_2)
		}
		BALL, PLAYER_1_SCORE, PLAYER_2_SCORE, HAS_SCORED =
			CollideBoundaryBall(WINDOW, BALL, PLAYER_1_SCORE, PLAYER_2_SCORE)
		BALL = ApplyVelocityBall(BALL, FRAME_DURATION)
		BALL = UpdatePreviousPosition(BALL)
		if HAS_SCORED {
			fmt.Printf("score! p1: %d, p2: %d\n", PLAYER_1_SCORE, PLAYER_2_SCORE)
		}
		BALL = HandleGameReset(HAS_SCORED, BALL, WINDOW)

		// rendering
		if HAS_SCORED {
			ClearWindow(Color_s{30, 30, 30, 255}, SF_WINDOW)
		} else {
			ClearWindow(Color_s{0, 0, 0, 255}, SF_WINDOW)
		}
		RenderRectangle(PLAYER_1.Rectangle, SF_WINDOW)
		RenderRectangle(PLAYER_2.Rectangle, SF_WINDOW)
		RenderRectangle(BALL.Rectangle, SF_WINDOW)
		RenderCenterLine(WINDOW, SF_WINDOW)
		DisplayWindow(SF_WINDOW)

		// handle timing
		CURRENT_FRAME_DURATION = UpdateCurrentFrameDuration(CURRENT_FRAME_START_TIME)
		Sleep(FRAME_DURATION, CURRENT_FRAME_DURATION)
	}
}
