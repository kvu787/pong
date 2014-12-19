package main

import (
	sf "bitbucket.org/kvu787/gosfml2"
	"time"
)

var (
	WINDOW                   Rectangle_s
	SF_WINDOW                *sf.RenderWindow
	FPS                      float64
	FRAME_DURATION           time.Duration
	CURRENT_FRAME_START_TIME time.Time
	CURRENT_FRAME_DURATION   time.Duration
	PADDLE_OFFSET            float64
	PLAYER_1                 Paddle_s
	PLAYER_2                 Paddle_s
	BALL                     Ball_s
	INPUT                    Input_s
	PLAYER_1_SCORE           int
	PLAYER_2_SCORE           int
	HAS_SCORED               bool
)
