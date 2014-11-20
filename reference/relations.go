package main

import (
	"os"
	"time"

	sf "bitbucket.org/kvu787/gosfml2"
)

func UpdateInput(sfWindow *sf.RenderWindow, input Input_s) Input_s {
	var event sf.Event
	for event = sfWindow.PollEvent(); event != nil; event = sfWindow.PollEvent() {
		switch e := event.(type) {
		case sf.EventKeyPressed:
			switch e.Code {
			case sf.KeyDown:
				input.IsDownArrowClicked = true
			case sf.KeyUp:
				input.IsUpArrowClicked = true
			}
		case sf.EventKeyReleased:
			switch e.Code {
			case sf.KeyDown:
				input.IsDownArrowClicked = false
			case sf.KeyUp:
				input.IsUpArrowClicked = false
			}
		case sf.EventClosed:
			input.IsWindowClosed = true
			sfWindow.Close()
		}
	}
	return input
}

func MovePlayer(input Input_s, duration time.Duration, player Paddle_s) Paddle_s {
	var up bool = input.IsUpArrowClicked
	var down bool = input.IsDownArrowClicked
	var position Vector_s = player.Rectangle.Position

	if up && !down {
		var displacement Vector_s = Vector_s{X: 0, Y: -player.Speed * duration.Seconds()}
		player.Rectangle.Position = VectorAdd(position, displacement)
	} else if down && !up {
		var displacement Vector_s = Vector_s{X: 0, Y: player.Speed * duration.Seconds()}
		player.Rectangle.Position = VectorAdd(position, displacement)
	}

	return player
}

func HandleWindowClose(input Input_s) {
	if input.IsWindowClosed {
		os.Exit(0)
	}
}

func UpdateCurrentFrameStartTime() time.Time {
	return time.Now()
}

func UpdateCurrentFrameDuration(startTime time.Time) time.Duration {
	return time.Since(startTime)
}

func Sync(frameDuration time.Duration, currentFrameDuration time.Duration) {
	var sleepDuration time.Duration = frameDuration - currentFrameDuration
	time.Sleep(sleepDuration)
}

func CollidePlayerBall(player Paddle_s, ball Ball_s, offset float64) Ball_s {
	var paddleSides [4]Segment_s = RectangleSegments(player.Rectangle)
	var i int
	for i = 0; i < 4; i++ {
		var paddleSide Segment_s = paddleSides[i]
		if AreCircleSegmentIntersecting(ball.Circle, paddleSide) {
			ball.Velocity = VectorMul(-1, ball.Velocity)

			// separate the segment and the circle
			var perpendicular Vector_s = PerpendicularVectorFromSegmentToPoint(paddleSide, ball.Circle.Position)
			var separation float64 = ball.Circle.Radius - VectorMagnitude(perpendicular)
			ball.Circle.Position = VectorAdd(ball.Circle.Position, VectorScale(separation+offset, perpendicular))
		}
	}
	return ball
}

func CollideBoundaryBall(window Rectangle_s, ball Ball_s) Ball_s {
	var boundarySides [4]Segment_s = RectangleSegments(window)
	var i int
	for i = 0; i < 4; i++ {
		var boundarySide Segment_s = boundarySides[i]
		if AreCircleSegmentIntersecting(ball.Circle, boundarySide) {
			var boundaryVector = VectorSub(boundarySide.End, boundarySide.Start)
			ball.Velocity = Reflect(ball.Velocity, boundaryVector)
		}
	}
	return ball
}

func ApplyVelocityBall(ball Ball_s, duration time.Duration) Ball_s {
	var displacement Vector_s = VectorMul(duration.Seconds(), ball.Velocity)
	ball.Circle.Position = VectorAdd(displacement, ball.Circle.Position)
	return ball
}

func MoveAi(ai Paddle_s, ball Ball_s, duration time.Duration) Paddle_s {
	var fromAiToBallY float64 = VectorSub(ball.Circle.Position, ai.Rectangle.Position).Y
	if fromAiToBallY < 0 {
		var displacement Vector_s = Vector_s{X: 0, Y: -ai.Speed * duration.Seconds()}
		ai.Rectangle.Position = VectorAdd(ai.Rectangle.Position, displacement)
	} else if fromAiToBallY > 0 {
		var displacement Vector_s = Vector_s{X: 0, Y: ai.Speed * duration.Seconds()}
		ai.Rectangle.Position = VectorAdd(ai.Rectangle.Position, displacement)
	}
	return ai
}
