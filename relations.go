package main

import (
	"math"
	"math/rand"
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
			case sf.KeyW:
				input.IsWClicked = true
			case sf.KeyS:
				input.IsSClicked = true
			}
		case sf.EventKeyReleased:
			switch e.Code {
			case sf.KeyDown:
				input.IsDownArrowClicked = false
			case sf.KeyUp:
				input.IsUpArrowClicked = false
			case sf.KeyW:
				input.IsWClicked = false
			case sf.KeyS:
				input.IsSClicked = false
			}
		case sf.EventClosed:
			input.IsWindowClosed = true
			sfWindow.Close()
		}
	}
	return input
}

func MovePlayer(input Input_s, duration time.Duration, player Paddle_s, playerNumber int) Paddle_s {
	var up bool
	var down bool

	if playerNumber == 1 {
		up = input.IsWClicked
		down = input.IsSClicked
	} else if playerNumber == 2 {
		up = input.IsUpArrowClicked
		down = input.IsDownArrowClicked
	} else {
		panic("MovePlayer: invalid player number")
	}

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

func KeepPlayerInBoundary(player Paddle_s, boundary Rectangle_s) Paddle_s {
	var minY float64 = player.Rectangle.Height / 2.0
	var maxY float64 = boundary.Height - player.Rectangle.Height/2.0
	player.Rectangle.Position.Y = Clamp(minY, player.Rectangle.Position.Y, maxY)

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

func Sleep(frameDuration time.Duration, currentFrameDuration time.Duration) {
	var sleepDuration time.Duration = frameDuration - currentFrameDuration
	time.Sleep(sleepDuration)
}

func CollidePaddleBall(player Paddle_s, ball Ball_s, offset float64, angleVarianceRange float64) Ball_s {
	var paddleSides [4]Segment_s = RectangleSegments(player.Rectangle)
	var i int
	for i = 0; i < 4; i++ {
		var paddleSide Segment_s = paddleSides[i]
		if AreCircleSegmentIntersecting(ball.Circle, paddleSide) {
			if i%2 == 0 {
				// long sides

				var resultant Vector_s = something(paddleSide, ball.Circle.Position)
				resultant = VectorScale(VectorMagnitude(ball.Velocity), resultant)
				ball.Velocity = resultant
			} else {
				// short sides

				ball.Velocity = VectorMul(-1, ball.Velocity)

				// randomize velocity direction
				var offset float64 = (rand.Float64() * DegreesToRadians(160)) - DegreesToRadians(80)
				var rejection = PerpendicularVectorFromLineToPoint(paddleSide, ball.Circle.Position)
				rejection = VectorScale(VectorMagnitude(ball.Velocity), rejection)
				rejection = VectorRotate(rejection, offset)
				ball.Velocity = rejection
			}

			// separate the segment and the circle
			var perpendicular Vector_s = PerpendicularVectorFromLineToPoint(paddleSide, ball.Circle.Position)
			var separation float64 = ball.Circle.Radius - VectorMagnitude(perpendicular)
			ball.Circle.Position = VectorAdd(ball.Circle.Position, VectorScale(separation+offset, perpendicular))
		}
	}
	return ball
}

func something(s Segment_s, p Vector_s) Vector_s {
	var offsets [8]float64 = [8]float64{20, 40, 60, 80, 100, 120, 140, 160}
	var endpoint Vector_s = s.Start
	var segmentVector = VectorSub(s.Start, s.End)
	var b Vector_s = PerpendicularVectorFromLineToPoint(s, p)
	var pointOnLine Vector_s = VectorSub(p, b)
	var distanceFromEnd float64 = VectorMagnitude(VectorSub(pointOnLine, endpoint))
	var sectionLength float64 = SegmentLength(s) / float64(len(offsets))
	var offsetIndex int = Floor(distanceFromEnd / sectionLength)
	var xProduct float64 = NormalizeScalar(VectorCrossProduct(segmentVector, b))
	return VectorNormalize(VectorRotate(segmentVector, xProduct*DegreesToRadians(offsets[offsetIndex])))
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
	var maxDisplacement float64 = ai.Speed * duration.Seconds()
	var displacementFromBall float64 = math.Abs(ball.Circle.Position.Y - ai.Rectangle.Position.Y)
	var displacement float64 = math.Min(maxDisplacement, displacementFromBall)
	if fromAiToBallY < 0 {
		ai.Rectangle.Position.Y -= displacement
	} else if fromAiToBallY > 0 {
		ai.Rectangle.Position.Y += displacement
	}
	return ai
}
