package ui

import (
	"fmt"
	"image"
	"log"
	"time"

	"github.com/giongto35/game-online/nes"
)

type View interface {
	Enter()
	Exit()
	Update(t, dt float64)
}

type Director struct {
	audio        *Audio
	view         View
	timestamp    float64
	imageChannel chan *image.RGBA
	inputChannel chan int
}

func NewDirector(audio *Audio, imageChannel chan *image.RGBA, inputChannel chan int) *Director {
	director := Director{}
	director.audio = audio
	director.imageChannel = imageChannel
	director.inputChannel = inputChannel
	return &director
}

func (d *Director) SetView(view View) {
	if d.view != nil {
		d.view.Exit()
	}
	d.view = view
	if d.view != nil {
		d.view.Enter()
	}
	d.timestamp = float64(time.Now().Nanosecond()) / float64(time.Second)
}

func (d *Director) Step() {
	//gl.Clear(gl.COLOR_BUFFER_BIT)
	timestamp := float64(time.Now().Nanosecond()) / float64(time.Second)
	fmt.Println("Time stamp", timestamp)
	dt := timestamp - d.timestamp
	fmt.Println("dt", dt)
	d.timestamp = timestamp
	fmt.Println("view", d.view)
	if d.view != nil {
		d.view.Update(timestamp, dt)
	}
}

func (d *Director) Start(paths []string) {
	if len(paths) == 1 {
		d.PlayGame(paths[0])
	}
	d.Run()
}

func (d *Director) Run() {
	for {
		d.Step()
	}
	d.SetView(nil)
}

func (d *Director) PlayGame(path string) {
	hash, err := hashFile(path)
	if err != nil {
		log.Fatalln(err)
	}
	console, err := nes.NewConsole(path)
	if err != nil {
		log.Fatalln(err)
	}
	d.SetView(NewGameView(d, console, path, hash, d.imageChannel, d.inputChannel))
}