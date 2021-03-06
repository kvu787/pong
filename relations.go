package main

import (
	"fmt"
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

func CollidePaddleBall(player Paddle_s, ball Ball_s) (Ball_s, bool) {
	var paddleSides [4]Segment_s = RectangleSegments(player.Rectangle)
	var i int
	var hasHitPaddle bool = false
	for i = 0; i < 4; i++ {
		var paddleSide Segment_s = paddleSides[i]
		if AreRectangleSegmentIntersecting(ball.Rectangle, paddleSide) {
			hasHitPaddle = true
			if i%2 == 0 {
				// long sides

				// reflect angle depends on what part of paddle is hit
				var resultant Vector_s = something(paddleSide, ball.Rectangle.Position)
				resultant = VectorScale(VectorMagnitude(ball.Velocity), resultant)
				ball.Velocity = resultant
			} else {
				// short sides

				// just reflect the ball
				ball.Velocity = VectorMul(-1, ball.Velocity)
			}

			// separate the segment and the ball
			// move the ball to its previous position
			ball.Rectangle.Position = ball.PreviousPosition

			// accelerate ball
			ball.Velocity = NewPolar(
				Clamp(0, VectorMagnitude(ball.Velocity)+BALL_ACCELERATION, BALL_MAX_VELOCITY),
				VectorAngle(ball.Velocity))

			// only collide with one side
			break
		}
	}
	return ball, hasHitPaddle
}

func something(s Segment_s, p Vector_s) Vector_s {
	var offsets [8]float64 = [8]float64{45, 57.86, 70.71, 83.57, 96.43, 109.29, 122.14, 135}
	var endpoint Vector_s = s.Start
	var segmentVector = VectorSub(s.Start, s.End)
	var b Vector_s = PerpendicularVectorFromLineToPoint(s, p)
	var pointOnLine Vector_s = VectorSub(p, b)
	var distanceFromEnd float64 = VectorMagnitude(VectorSub(pointOnLine, endpoint))
	var sectionLength float64 = SegmentLength(s) / float64(len(offsets))
	var offsetIndex int = Floor(distanceFromEnd / sectionLength)
	if offsetIndex >= len(offsets) {
		offsetIndex = len(offsets) - 1
	}
	var xProduct float64 = NormalizeScalar(VectorCrossProduct(segmentVector, b))
	return VectorNormalize(VectorRotate(segmentVector, xProduct*DegreesToRadians(offsets[offsetIndex])))
}

func CollideBoundaryBall(window Rectangle_s, ball Ball_s, p1Score int, p2Score int) (Ball_s, int, int, bool) {
	var boundarySides [4]Segment_s = RectangleSegments(window)
	var hasP1Scored bool = false
	var hasP2Scored bool = false
	var i int
	for i = 0; i < 4; i++ {
		var boundarySide Segment_s = boundarySides[i]
		if AreRectangleSegmentIntersecting(ball.Rectangle, boundarySide) {

			// reflect ball on boundary
			var boundaryVector = VectorSub(boundarySide.End, boundarySide.Start)
			ball.Velocity = Reflect(ball.Velocity, boundaryVector)
			ball.Rectangle.Position = ball.PreviousPosition

			// separate
			ball.Rectangle.Position = ball.PreviousPosition

			// handle scoring
			if i == 0 {
				hasP1Scored = true
			} else if i == 2 {
				hasP2Scored = true
			}

			// only collide with one side
			break
		}
	}
	if hasP1Scored {
		return ball, p1Score + 1, p2Score, true
	} else if hasP2Scored {
		return ball, p1Score, p2Score + 1, true
	} else {
		return ball, p1Score, p2Score, false
	}
}

func ApplyVelocityBall(ball Ball_s, duration time.Duration) Ball_s {
	var displacement Vector_s = VectorMul(duration.Seconds(), ball.Velocity)
	ball.Rectangle.Position = VectorAdd(displacement, ball.Rectangle.Position)
	return ball
}

func UpdatePreviousPosition(ball Ball_s) Ball_s {
	ball.PreviousPosition = ball.Rectangle.Position
	return ball
}

func HandleGameReset(hasScored bool, ball Ball_s, window Rectangle_s) Ball_s {
	if hasScored {
		return Ball_s{
			Rectangle: Rectangle_s{
				Position: Vector_s{
					X: WINDOW.Width / 2,
					Y: WINDOW.Height / 2},
				Width:  10,
				Height: 10},
			Velocity: NewPolar(BALL_START_VELOCITY, GenerateRandomBallDirection())}
	} else {
		return ball
	}
}

func RenderCenterLine(window Rectangle_s, sfWindow *sf.RenderWindow) {
	var r Rectangle_s = Rectangle_s{
		Position: Vector_s{
			X: WINDOW.Width / 2,
			Y: WINDOW.Height / 2},
		Width:  3,
		Height: WINDOW.Height}
	RenderRectangle(r, sfWindow)
}

func ClampBall(ball Ball_s, player1 Paddle_s, player2 Paddle_s) Ball_s {
	var padding float64 = player1.Rectangle.Width/2 + ball.Rectangle.Width/2 + 2
	ball.Rectangle.Position.X = Clamp(
		player1.Rectangle.Position.X+padding,
		ball.Rectangle.Position.X,
		player2.Rectangle.Position.X-padding)
	return ball
}

func RenderScores(player1Score int, player2Score int, window Rectangle_s, sfWindow *sf.RenderWindow) {
	f, err := sf.NewFontFromFile("DroidSansMono.ttf")
	if err != nil {
		panic("could not load font")
	}

	t, err := sf.NewText(f)
	if err != nil {
		panic("could not create text object")
	}
	t.SetString(fmt.Sprintf("%d", player1Score))
	t.SetCharacterSize(50)
	t.SetPosition(sf.Vector2f{13, 10})
	t.SetColor(sf.Color{255, 255, 255, 255})
	SF_WINDOW.Draw(t, sf.DefaultRenderStates())

	t, err = sf.NewText(f)
	if err != nil {
		panic("could not create text object")
	}
	t.SetString(fmt.Sprintf("%d", player2Score))
	t.SetCharacterSize(50)
	t.SetPosition(sf.Vector2f{float32(window.Width/2.0 + 15), 10})
	t.SetColor(sf.Color{255, 255, 255, 255})
	SF_WINDOW.Draw(t, sf.DefaultRenderStates())
}
