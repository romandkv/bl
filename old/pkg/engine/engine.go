package engine

import (
	"fmt"
	"github.com/faiface/pixel"
	"github.com/faiface/pixel/pixelgl"
	"github.com/gammazero/deque"
	models2 "gitlab.equifax.local/dro/bl/old/pkg/models"
	"golang.org/x/image/colornames"
)

type Engine struct {
	Characters map[string]models2.Character
	Sounds map[string]models2.Sound
	Backgrounds map[string]models2.Background
	Sprites map[string]models2.Sprite
	Files []string
	Queue deque.Deque
	Window *pixelgl.Window
	LoadNextChapter chan bool
	LoadedNextChapter chan bool
	EndOfGame chan bool
	Cache bool
	CachedLabels map[string]deque.Deque
	CacheLableName string
}

func newEngine() Engine {
	return Engine{
		Characters:  make(map[string]models2.Character),
		Sounds:      make(map[string]models2.Sound),
		Backgrounds: make(map[string]models2.Background),
		Sprites:     make(map[string]models2.Sprite),
		LoadNextChapter : make(chan bool),
		LoadedNextChapter : make(chan bool),
		EndOfGame : make(chan bool),
		Cache: false,
		CachedLabels: make(map[string]deque.Deque),
	}
}

func Run() {
	var err error
	cfg := pixelgl.WindowConfig{
		Title:  "BL!",
		Bounds: pixel.R(0, 0, 1024, 768),
		VSync: true,
	}
	VNE.Window, err = pixelgl.NewWindow(cfg)
	if err != nil {
		panic(err)
	}
	VNE.Window.Clear(colornames.Skyblue)
	for !VNE.Window.Closed() {
		if VNE.Window.JustPressed(pixelgl.KeySpace) {
			fmt.Println("PRESSED SPACE")
			if VNE.Queue.Len() != 0 {
				VNE.Queue.PopFront().(Action).Execute()
			}
			fmt.Println("Очередь: " + fmt.Sprint(VNE.Queue.Len()))
			if VNE.Queue.Len() == 0 {
				select {
					case <- VNE.EndOfGame :
						return
					default:
						break
				}
				VNE.LoadNextChapter <- true
				<- VNE.LoadedNextChapter
			}
		}

		VNE.Window.Update()
	}
}

