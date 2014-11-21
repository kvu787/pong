package main

import (
	sf "bitbucket.org/kvu787/gosfml2"
	"os"
)

func UpdateInput(sfWindow *sf.RenderWindow) Input_s {
	var event sf.Event
	for event = sfWindow.PollEvent(); event != nil; event = sfWindow.PollEvent() {
		switch event.(type) {
		case sf.EventClosed:
			return Input_s{IsWindowClosed: true}
		}
	}
	return Input_s{IsWindowClosed: false}
}

func HandleWindowClose(input Input_s) {
	if input.IsWindowClosed {
		os.Exit(0)
	}
}

func SetupWindow(width int, height int) *sf.RenderWindow {
	return sf.NewRenderWindow(
		sf.VideoMode{Width: uint(width), Height: uint(height), BitsPerPixel: 32},
		"pong",
		sf.StyleDefault,
		sf.DefaultContextSettings())
}

func ClearWindow(window *sf.RenderWindow) {
	window.Clear(sf.ColorBlack())
}

func DisplayWindow(window *sf.RenderWindow) {
	window.Display()
}
